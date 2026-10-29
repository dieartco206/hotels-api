package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"hotels-api/config"
	"hotels-api/models"

	"github.com/gin-gonic/gin"
)

type HotelCreateInput struct {
	Name        string   `json:"name" binding:"required"`
	Tagline     string   `json:"tagline"`
	Type        string   `json:"type" binding:"required"`
	City        string   `json:"city" binding:"required"`
	Area        string   `json:"area"`
	Country     string   `json:"country" binding:"required"`
	Address     string   `json:"address" binding:"required"`
	StarRating  int      `json:"star_rating"`
	PriceStart  float64  `json:"price_start" binding:"required"`
	Currency    string   `json:"currency"`
	Image       string   `json:"image" binding:"required"`
	Gallery     string   `json:"gallery"`
	Facilities  string   `json:"facilities"`
	Description string   `json:"description"`
}

func GetHotels(c *gin.Context) {
	var hotels []models.Hotel
	query := config.DB.Model(&models.Hotel{}).Preload("Rooms")

	// Filter by status (public view only sees active)
	showAll := c.Query("all") == "true"
	if !showAll {
		query = query.Where("status = ?", models.HotelStatusActive)
	}

	// Filter by country
	country := strings.ToUpper(strings.TrimSpace(c.Query("country")))
	if country != "" && country != "ALL" {
		query = query.Where("country = ?", country)
	}

	// Filter by city
	city := strings.TrimSpace(c.Query("city"))
	if city != "" {
		query = query.Where("city LIKE ?", "%"+city+"%")
	}

	// Filter by type
	hotelType := strings.TrimSpace(c.Query("type"))
	if hotelType != "" && hotelType != "all" {
		query = query.Where("type = ?", hotelType)
	}

	// Keyword search
	search := strings.TrimSpace(c.Query("search"))
	if search != "" {
		term := "%" + search + "%"
		query = query.Where("name LIKE ? OR city LIKE ? OR area LIKE ? OR description LIKE ?", term, term, term, term)
	}

	// Star rating
	if minStarStr := c.Query("star_rating"); minStarStr != "" {
		if minStar, err := strconv.Atoi(minStarStr); err == nil && minStar > 0 {
			query = query.Where("star_rating >= ?", minStar)
		}
	}

	// Sorting
	sortBy := c.Query("sort")
	switch sortBy {
	case "price_asc":
		query = query.Order("price_start ASC")
	case "price_desc":
		query = query.Order("price_start DESC")
	case "rating_desc":
		query = query.Order("rating DESC, review_count DESC")
	default:
		query = query.Order("rating DESC, id ASC")
	}

	if err := query.Find(&hotels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch hotels"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(hotels),
		"data":    hotels,
	})
}

func GetHotelByID(c *gin.Context) {
	idParam := c.Param("id")
	hotelID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid hotel ID"})
		return
	}

	var hotel models.Hotel
	if err := config.DB.Preload("Rooms").Preload("Partner").First(&hotel, hotelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Hotel not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    hotel,
	})
}

func CreateHotel(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid user context"})
		return
	}

	var input HotelCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	currency := strings.ToUpper(input.Currency)
	if currency == "" {
		if input.Country == "SG" {
			currency = "SGD"
		} else if input.Country == "MY" {
			currency = "MYR"
		} else {
			currency = "IDR"
		}
	}

	star := input.StarRating
	if star <= 0 {
		star = 5
	}

	hotel := models.Hotel{
		Name:        strings.TrimSpace(input.Name),
		Tagline:     strings.TrimSpace(input.Tagline),
		Type:        strings.ToLower(input.Type),
		City:        strings.TrimSpace(input.City),
		Area:        strings.TrimSpace(input.Area),
		Country:     strings.ToUpper(input.Country),
		Address:     strings.TrimSpace(input.Address),
		StarRating:  star,
		Rating:      4.9,
		ReviewCount: 1,
		PriceStart:  input.PriceStart,
		Currency:    currency,
		Image:       input.Image,
		Gallery:     input.Gallery,
		Facilities:  input.Facilities,
		Description: input.Description,
		PartnerID:   userID,
		Status:      models.HotelStatusActive,
	}

	if err := config.DB.Create(&hotel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create property"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Hotel property created successfully",
		"data":    hotel,
	})
}

func UpdateHotel(c *gin.Context) {
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

	// Security: verify partner ownership unless admin
	roleVal, _ := c.Get("userRole")
	userRole, _ := roleVal.(models.UserRole)
	userIDVal, _ := c.Get("userID")
	userID, _ := userIDVal.(uint)

	if userRole != models.RoleAdmin && hotel.PartnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Not authorized to modify this property"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Strip un-editable fields
	delete(updateData, "id")
	delete(updateData, "partner_id")
	delete(updateData, "created_at")

	if err := config.DB.Model(&hotel).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update hotel"})
		return
	}

	config.DB.Preload("Rooms").First(&hotel, hotelID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hotel updated successfully",
		"data":    hotel,
	})
}
