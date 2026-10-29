package controllers

import (
	"net/http"
	"strconv"

	"hotels-api/config"
	"hotels-api/models"

	"github.com/gin-gonic/gin"
)

type RoomInput struct {
	HotelID          uint    `json:"hotel_id" binding:"required"`
	Name             string  `json:"name" binding:"required"`
	Type             string  `json:"type" binding:"required"`
	Price            float64 `json:"price" binding:"required"`
	OriginalPrice    float64 `json:"original_price"`
	MaxGuests        int     `json:"max_guests"`
	BedType          string  `json:"bed_type"`
	SizeSqm          int     `json:"size_sqm"`
	RemainingRooms   int     `json:"remaining_rooms"`
	TotalRooms       int     `json:"total_rooms"`
	FreeBreakfast    bool    `json:"free_breakfast"`
	FreeCancellation bool    `json:"free_cancellation"`
	Amenities        string  `json:"amenities"`
	Image            string  `json:"image"`
}

type AllotmentUpdateInput struct {
	RemainingRooms *int     `json:"remaining_rooms"`
	Price          *float64 `json:"price"`
	Status         *string  `json:"status"`
}

func CreateRoom(c *gin.Context) {
	var input RoomInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	var hotel models.Hotel
	if err := config.DB.First(&hotel, input.HotelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Hotel not found"})
		return
	}

	maxGuests := input.MaxGuests
	if maxGuests <= 0 {
		maxGuests = 2
	}
	remaining := input.RemainingRooms
	if remaining <= 0 {
		remaining = 5
	}
	total := input.TotalRooms
	if total <= 0 {
		total = 10
	}
	bed := input.BedType
	if bed == "" {
		bed = "1 King Bed"
	}

	room := models.Room{
		HotelID:          input.HotelID,
		Name:             input.Name,
		Type:             input.Type,
		Price:            input.Price,
		OriginalPrice:    input.OriginalPrice,
		MaxGuests:        maxGuests,
		BedType:          bed,
		SizeSqm:          input.SizeSqm,
		RemainingRooms:   remaining,
		TotalRooms:       total,
		FreeBreakfast:    input.FreeBreakfast,
		FreeCancellation: input.FreeCancellation,
		Amenities:        input.Amenities,
		Image:            input.Image,
		Status:           "available",
	}

	if err := config.DB.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create room"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Room created successfully",
		"data":    room,
	})
}

func UpdateRoom(c *gin.Context) {
	idParam := c.Param("id")
	roomID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid room ID"})
		return
	}

	var room models.Room
	if err := config.DB.First(&room, roomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Room not found"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	delete(updateData, "id")
	delete(updateData, "hotel_id")

	if err := config.DB.Model(&room).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update room"})
		return
	}

	config.DB.First(&room, roomID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Room updated successfully",
		"data":    room,
	})
}

func UpdateAllotment(c *gin.Context) {
	idParam := c.Param("id")
	roomID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid room ID"})
		return
	}

	var room models.Room
	if err := config.DB.First(&room, roomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Room not found"})
		return
	}

	var input AllotmentUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if input.RemainingRooms != nil {
		updates["remaining_rooms"] = *input.RemainingRooms
		if *input.RemainingRooms <= 0 {
			updates["status"] = "sold_out"
		} else {
			updates["status"] = "available"
		}
	}
	if input.Price != nil && *input.Price > 0 {
		updates["price"] = *input.Price
	}
	if input.Status != nil && *input.Status != "" {
		updates["status"] = *input.Status
	}

	if len(updates) > 0 {
		if err := config.DB.Model(&room).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update allotment"})
			return
		}
	}

	config.DB.First(&room, roomID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Room inventory allotment updated successfully",
		"data":    room,
	})
}
