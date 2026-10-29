package controllers

import (
	"net/http"
	"strconv"
	"time"

	"hotels-api/config"
	"hotels-api/models"

	"github.com/gin-gonic/gin"
)

func GetAdminDashboard(c *gin.Context) {
	var totalHotels int64
	var activeHotels int64
	var totalBookings int64
	var totalPartners int64
	var pendingApplications int64

	config.DB.Model(&models.Hotel{}).Count(&totalHotels)
	config.DB.Model(&models.Hotel{}).Where("status = ?", models.HotelStatusActive).Count(&activeHotels)
	config.DB.Model(&models.Booking{}).Count(&totalBookings)
	config.DB.Model(&models.User{}).Where("role = ?", models.RolePartner).Count(&totalPartners)
	config.DB.Model(&models.PartnerApplication{}).Where("status = ?", "pending").Count(&pendingApplications)

	// Financials
	var totalGBV float64
	var totalCommission float64
	var pendingPayouts float64

	config.DB.Model(&models.Settlement{}).Where("status != ?", models.SettlementRejected).Select("COALESCE(SUM(gross_amount), 0)").Scan(&totalGBV)
	config.DB.Model(&models.Settlement{}).Where("status != ?", models.SettlementRejected).Select("COALESCE(SUM(platform_commission), 0)").Scan(&totalCommission)
	config.DB.Model(&models.Settlement{}).Where("status = ?", models.SettlementPending).Select("COALESCE(SUM(net_payout), 0)").Scan(&pendingPayouts)

	// Country breakdown
	type CountryStat struct {
		Country string `json:"country"`
		Count   int64  `json:"count"`
	}
	var countryStats []CountryStat
	config.DB.Model(&models.Hotel{}).Select("country, count(*) as count").Group("country").Scan(&countryStats)

	// Recent Bookings across platform
	var recentBookings []models.Booking
	config.DB.Preload("Hotel").Preload("Room").Order("id DESC").Limit(10).Find(&recentBookings)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_gbv":            totalGBV,
			"net_revenue":          totalCommission,
			"pending_payouts":      pendingPayouts,
			"total_hotels":         totalHotels,
			"active_hotels":        activeHotels,
			"total_bookings":       totalBookings,
			"total_partners":       totalPartners,
			"pending_applications": pendingApplications,
			"country_distribution": countryStats,
			"recent_bookings":      recentBookings,
		},
	})
}

func GetPartnerApplications(c *gin.Context) {
	var apps []models.PartnerApplication
	query := config.DB.Order("id DESC")

	if status := c.Query("status"); status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&apps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch partner applications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(apps),
		"data":    apps,
	})
}

func ReviewPartnerApplication(c *gin.Context) {
	idParam := c.Param("id")
	appID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid application ID"})
		return
	}

	var app models.PartnerApplication
	if err := config.DB.First(&app, appID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Application not found"})
		return
	}

	var input struct {
		Action string `json:"action" binding:"required"` // approve, reject
		Notes  string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if input.Action == "approve" {
		app.Status = "approved"
	} else if input.Action == "reject" {
		app.Status = "rejected"
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid action. Use 'approve' or 'reject'"})
		return
	}

	if input.Notes != "" {
		app.Notes = input.Notes
	}

	if err := config.DB.Save(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Application status updated to " + app.Status,
		"data":    app,
	})
}

func GetAllBookings(c *gin.Context) {
	var bookings []models.Booking
	query := config.DB.Preload("Hotel").Preload("Room").Order("id DESC")

	if status := c.Query("status"); status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}
	if country := c.Query("country"); country != "" && country != "all" {
		query = query.Joins("JOIN hotels ON hotels.id = bookings.hotel_id").Where("hotels.country = ?", country)
	}

	if err := query.Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(bookings),
		"data":    bookings,
	})
}

func GetAllSettlements(c *gin.Context) {
	var settlements []models.Settlement
	query := config.DB.Preload("Hotel").Preload("Booking").Order("id DESC")

	if status := c.Query("status"); status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&settlements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch settlements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(settlements),
		"data":    settlements,
	})
}

func ApproveSettlementPayout(c *gin.Context) {
	idParam := c.Param("id")
	settlementID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid settlement ID"})
		return
	}

	var settlement models.Settlement
	if err := config.DB.First(&settlement, settlementID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Settlement record not found"})
		return
	}

	now := time.Now()
	settlement.Status = models.SettlementTransferred
	settlement.PayoutDate = &now

	if err := config.DB.Save(&settlement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to disburse settlement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Settlement payout approved and transferred successfully",
		"data":    settlement,
	})
}

func ToggleHotelStatus(c *gin.Context) {
	idParam := c.Param("id")
	hotelID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid hotel ID"})
		return
	}

	var hotel models.Hotel
	if err := config.DB.First(&hotel, hotelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Hotel not found"})
		return
	}

	var input struct {
		Status models.HotelStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	hotel.Status = input.Status
	if err := config.DB.Save(&hotel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update property status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Property status updated to " + string(hotel.Status),
		"data":    hotel,
	})
}
