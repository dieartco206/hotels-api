package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"hotels-api/config"
	"hotels-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingCreateInput struct {
	HotelID         uint      `json:"hotel_id" binding:"required"`
	RoomID          uint      `json:"room_id" binding:"required"`
	GuestName       string    `json:"guest_name" binding:"required"`
	GuestEmail      string    `json:"guest_email" binding:"required,email"`
	GuestPhone      string    `json:"guest_phone" binding:"required"`
	CheckIn         string    `json:"check_in" binding:"required"`  // YYYY-MM-DD
	CheckOut        string    `json:"check_out" binding:"required"` // YYYY-MM-DD
	GuestsCount     int       `json:"guests_count"`
	PaymentMethod   string    `json:"payment_method"`
	SpecialRequests string    `json:"special_requests"`
}

func generateBookingCode(country string) string {
	if country == "" {
		country = "SG"
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(90000) + 10000
	return fmt.Sprintf("HTL-%s-%d", country, num)
}

func CreateBooking(c *gin.Context) {
	var input BookingCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	checkInDate, err := time.Parse("2006-01-02", input.CheckIn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid check_in date format. Use YYYY-MM-DD"})
		return
	}

	checkOutDate, err := time.Parse("2006-01-02", input.CheckOut)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid check_out date format. Use YYYY-MM-DD"})
		return
	}

	if !checkOutDate.After(checkInDate) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Check-out date must be after check-in date"})
		return
	}

	nights := int(checkOutDate.Sub(checkInDate).Hours() / 24)
	if nights <= 0 {
		nights = 1
	}

	var hotel models.Hotel
	if err := config.DB.First(&hotel, input.HotelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Hotel not found"})
		return
	}

	var room models.Room
	if err := config.DB.First(&room, input.RoomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Room not found"})
		return
	}

	if room.RemainingRooms <= 0 {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Selected room category is sold out for these dates"})
		return
	}

	var userIDPtr *uint
	if uidVal, exists := c.Get("userID"); exists {
		if uid, ok := uidVal.(uint); ok && uid > 0 {
			userIDPtr = &uid
		}
	}

	totalAmount := room.Price * float64(nights)
	bookingCode := generateBookingCode(hotel.Country)
	qrData := fmt.Sprintf("HTL_CHECKIN:%s:%s", bookingCode, strings.ToUpper(strings.TrimSpace(input.GuestName)))

	paymentMethod := input.PaymentMethod
	if paymentMethod == "" {
		paymentMethod = "credit_card"
	}

	guestsCount := input.GuestsCount
	if guestsCount <= 0 {
		guestsCount = 1
	}

	booking := models.Booking{
		BookingCode:     bookingCode,
		HotelID:         hotel.ID,
		RoomID:          room.ID,
		UserID:          userIDPtr,
		GuestName:       strings.TrimSpace(input.GuestName),
		GuestEmail:      strings.ToLower(strings.TrimSpace(input.GuestEmail)),
		GuestPhone:      strings.TrimSpace(input.GuestPhone),
		CheckIn:         checkInDate,
		CheckOut:        checkOutDate,
		Nights:          nights,
		GuestsCount:     guestsCount,
		TotalAmount:     totalAmount,
		Currency:        hotel.Currency,
		Status:          models.BookingConfirmed,
		PaymentMethod:   paymentMethod,
		PaymentStatus:   "paid",
		SpecialRequests: strings.TrimSpace(input.SpecialRequests),
		QrCodeData:      qrData,
	}

	// Calculate commission (15%)
	platformCommission := totalAmount * 0.15
	netPayout := totalAmount - platformCommission

	settlementCode := fmt.Sprintf("STL-%s-%d", hotel.Country, rand.Intn(900000)+100000)

	settlement := models.Settlement{
		SettlementCode:     settlementCode,
		PartnerID:          hotel.PartnerID,
		HotelID:            hotel.ID,
		GrossAmount:        totalAmount,
		PlatformCommission: platformCommission,
		NetPayout:          netPayout,
		Currency:           hotel.Currency,
		Status:             models.SettlementPending,
		BankName:           "DBS Singapore / BCA Indo / Maybank",
		BankAccount:        "8892-0192-4410",
	}

	// Transaction to save booking, decrement room allotment, and create settlement
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		// Decrement room inventory
		newRemaining := room.RemainingRooms - 1
		statusUpdate := room.Status
		if newRemaining <= 0 {
			newRemaining = 0
			statusUpdate = "sold_out"
		}
		if err := tx.Model(&room).Updates(map[string]interface{}{
			"remaining_rooms": newRemaining,
			"status":          statusUpdate,
		}).Error; err != nil {
			return err
		}

		// Link settlement with booking ID
		settlement.BookingID = booking.ID
		if err := tx.Create(&settlement).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to complete reservation transaction: " + err.Error()})
		return
	}

	config.DB.Preload("Hotel").Preload("Room").First(&booking, booking.ID)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Reservation confirmed successfully",
		"data":    booking,
	})
}

func GetBookingByCode(c *gin.Context) {
	code := strings.TrimSpace(c.Param("code"))
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Booking code is required"})
		return
	}

	var booking models.Booking
	if err := config.DB.Preload("Hotel").Preload("Room").Where("booking_code = ?", code).First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Booking not found with code: " + code})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    booking,
	})
}

func GetMyBookings(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Invalid user ID"})
		return
	}

	var bookings []models.Booking
	if err := config.DB.Preload("Hotel").Preload("Room").Where("user_id = ?", userID).Order("id DESC").Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(bookings),
		"data":    bookings,
	})
}

func CheckInGuest(c *gin.Context) {
	var input struct {
		CodeOrQR string `json:"code_or_qr" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	cleaned := strings.TrimSpace(input.CodeOrQR)
	var bookingCode string

	if strings.HasPrefix(cleaned, "HTL_CHECKIN:") {
		parts := strings.Split(cleaned, ":")
		if len(parts) >= 2 {
			bookingCode = parts[1]
		}
	} else {
		bookingCode = cleaned
	}

	var booking models.Booking
	if err := config.DB.Preload("Hotel").Preload("Room").Where("booking_code = ?", bookingCode).First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Booking code not found: " + bookingCode})
		return
	}

	if booking.Status == models.BookingCancelled {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Cannot check in: this reservation has been cancelled"})
		return
	}

	if booking.Status == models.BookingCheckedIn {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Guest is already checked in",
			"data":    booking,
		})
		return
	}

	now := time.Now()
	booking.Status = models.BookingCheckedIn
	booking.CheckedInAt = &now

	if err := config.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update check-in status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Guest checked in successfully! Key card issued.",
		"data":    booking,
	})
}

func CancelBooking(c *gin.Context) {
	idParam := c.Param("id")
	bookingID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid booking ID"})
		return
	}

	var booking models.Booking
	if err := config.DB.First(&booking, bookingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Booking not found"})
		return
	}

	if booking.Status == models.BookingCancelled {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Booking is already cancelled"})
		return
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		booking.Status = models.BookingCancelled
		booking.PaymentStatus = "refunded"
		if err := tx.Save(&booking).Error; err != nil {
			return err
		}

		// Restore room allotment
		var room models.Room
		if err := tx.First(&room, booking.RoomID).Error; err == nil {
			tx.Model(&room).Updates(map[string]interface{}{
				"remaining_rooms": room.RemainingRooms + 1,
				"status":          "available",
			})
		}

		// Void settlement
		tx.Model(&models.Settlement{}).Where("booking_id = ?", booking.ID).Update("status", models.SettlementRejected)

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to cancel reservation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Booking cancelled and room inventory restored",
	})
}
