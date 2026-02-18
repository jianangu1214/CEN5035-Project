package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/travellog/backend/database"
	"github.com/travellog/backend/models"
)

type HotelHandler struct{}

func NewHotelHandler() *HotelHandler {
	return &HotelHandler{}
}

func getUserID(c *gin.Context) (uint, bool) {
	id, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	uid, ok := id.(uint)
	return uid, ok
}

// GetAllHotels - GET /hotels
func (h *HotelHandler) GetAllHotels(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var hotels []models.Hotel
	if err := database.GetDB().Where("user_id = ?", userID).Find(&hotels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hotels"})
		return
	}

	var responses []models.HotelResponse
	for _, hotel := range hotels {
		responses = append(responses, hotel.ToResponse())
	}

	if responses == nil {
		responses = []models.HotelResponse{}
	}

	c.JSON(http.StatusOK, gin.H{"hotels": responses})
}

// GetHotel - GET /hotels/:id
func (h *HotelHandler) GetHotel(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	hotelID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID"})
		return
	}

	var hotel models.Hotel
	if err := database.GetDB().Where("id = ? AND user_id = ?", hotelID, userID).First(&hotel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hotel": hotel.ToResponse()})
}

// CreateHotel - POST /hotels
func (h *HotelHandler) CreateHotel(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.CreateHotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Parse dates
	checkIn, err := time.Parse("2006-01-02", req.CheckIn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check_in date format, use YYYY-MM-DD"})
		return
	}

	checkOut, err := time.Parse("2006-01-02", req.CheckOut)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check_out date format, use YYYY-MM-DD"})
		return
	}

	hotel := models.Hotel{
		UserID:    userID,
		HotelName: req.HotelName,
		City:      req.City,
		Country:   req.Country,
		CheckIn:   checkIn,
		CheckOut:  checkOut,
		Price:     req.Price,
		Notes:     req.Notes,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	if err := database.GetDB().Create(&hotel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create hotel"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"hotel": hotel.ToResponse()})
}

// UpdateHotel - PUT /hotels/:id
func (h *HotelHandler) UpdateHotel(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	hotelID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID"})
		return
	}

	var hotel models.Hotel
	if err := database.GetDB().Where("id = ? AND user_id = ?", hotelID, userID).First(&hotel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
		return
	}

	var req models.UpdateHotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update fields if provided
	if req.HotelName != "" {
		hotel.HotelName = req.HotelName
	}
	if req.City != "" {
		hotel.City = req.City
	}
	if req.Country != "" {
		hotel.Country = req.Country
	}
	if req.CheckIn != "" {
		checkIn, err := time.Parse("2006-01-02", req.CheckIn)
		if err == nil {
			hotel.CheckIn = checkIn
		}
	}
	if req.CheckOut != "" {
		checkOut, err := time.Parse("2006-01-02", req.CheckOut)
		if err == nil {
			hotel.CheckOut = checkOut
		}
	}
	if req.Price > 0 {
		hotel.Price = req.Price
	}

	if req.Notes != "" {
		hotel.Notes = req.Notes
	}
	if req.Latitude != 0 {
		hotel.Latitude = req.Latitude
	}
	if req.Longitude != 0 {
		hotel.Longitude = req.Longitude
	}

	if err := database.GetDB().Save(&hotel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hotel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hotel": hotel.ToResponse()})
}

// DeleteHotel - DELETE /hotels/:id
func (h *HotelHandler) DeleteHotel(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	hotelID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID"})
		return
	}

	var hotel models.Hotel
	if err := database.GetDB().Where("id = ? AND user_id = ?", hotelID, userID).First(&hotel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
		return
	}

	if err := database.GetDB().Delete(&hotel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete hotel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hotel deleted successfully"})
}