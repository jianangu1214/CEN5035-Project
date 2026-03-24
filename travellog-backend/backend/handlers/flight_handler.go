package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/travellog/backend/database"
	"github.com/travellog/backend/models"
)

type FlightHandler struct{}

func NewFlightHandler() *FlightHandler {
	return &FlightHandler{}
}

// getUserID retrieves the authenticated user id from context.
// Duplicated here to keep handlers decoupled; matches hotel handler usage.
func getUserIDFromContext(c *gin.Context) (uint, bool) {
	id, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	uid, ok := id.(uint)
	return uid, ok
}

// GetAllFlights - GET /flights
func (h *FlightHandler) GetAllFlights(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var flights []models.Flight
	if err := database.GetDB().Where("user_id = ?", userID).Find(&flights).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch flights"})
		return
	}

	responses := make([]models.FlightResponse, 0, len(flights))
	for _, f := range flights {
		responses = append(responses, f.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{"flights": responses})
}

// GetFlight - GET /flights/:id
func (h *FlightHandler) GetFlight(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	flightID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	var flight models.Flight
	if err := database.GetDB().Where("id = ? AND user_id = ?", flightID, userID).First(&flight).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"flight": flight.ToResponse()})
}

// CreateFlight - POST /flights
func (h *FlightHandler) CreateFlight(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.CreateFlightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	depart, err := time.Parse(time.RFC3339, req.DepartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid depart_time format, use RFC3339"})
		return
	}

	arrive, err := time.Parse(time.RFC3339, req.ArriveTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid arrive_time format, use RFC3339"})
		return
	}

	flight := models.Flight{
		UserID:       userID,
		Airline:      req.Airline,
		FlightNumber: req.FlightNumber,
		FromAirport:  req.FromAirport,
		ToAirport:    req.ToAirport,
		DepartTime:   depart,
		ArriveTime:   arrive,
		Price:        req.Price,
		Notes:        req.Notes,
	}

	if err := database.GetDB().Create(&flight).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create flight"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"flight": flight.ToResponse()})
}

// UpdateFlight - PUT /flights/:id
func (h *FlightHandler) UpdateFlight(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	flightID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	var flight models.Flight
	if err := database.GetDB().Where("id = ? AND user_id = ?", flightID, userID).First(&flight).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	var req models.UpdateFlightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Airline != "" {
		flight.Airline = req.Airline
	}
	if req.FlightNumber != "" {
		flight.FlightNumber = req.FlightNumber
	}
	if req.FromAirport != "" {
		flight.FromAirport = req.FromAirport
	}
	if req.ToAirport != "" {
		flight.ToAirport = req.ToAirport
	}
	if req.DepartTime != "" {
		if depart, err := time.Parse(time.RFC3339, req.DepartTime); err == nil {
			flight.DepartTime = depart
		}
	}
	if req.ArriveTime != "" {
		if arrive, err := time.Parse(time.RFC3339, req.ArriveTime); err == nil {
			flight.ArriveTime = arrive
		}
	}
	if req.Price > 0 {
		flight.Price = req.Price
	}
	if req.Notes != "" {
		flight.Notes = req.Notes
	}

	if err := database.GetDB().Save(&flight).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update flight"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"flight": flight.ToResponse()})
}

// DeleteFlight - DELETE /flights/:id
func (h *FlightHandler) DeleteFlight(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	flightID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	var flight models.Flight
	if err := database.GetDB().Where("id = ? AND user_id = ?", flightID, userID).First(&flight).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	if err := database.GetDB().Delete(&flight).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete flight"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flight deleted successfully"})
}
