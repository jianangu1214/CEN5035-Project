package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/travellog/backend/database"
	"github.com/travellog/backend/models"
)

type MapHandler struct{}

func NewMapHandler() *MapHandler {
	return &MapHandler{}
}

type mapHotel struct {
	ID        uint    `json:"id"`
	HotelName string  `json:"hotel_name"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type mapFlight struct {
	ID           uint    `json:"id"`
	Airline      string  `json:"airline"`
	FlightNumber string  `json:"flight_number"`
	FromAirport  string  `json:"from_airport"`
	ToAirport    string  `json:"to_airport"`
	FromLat      float64 `json:"from_lat"`
	FromLng      float64 `json:"from_lng"`
	ToLat        float64 `json:"to_lat"`
	ToLng        float64 `json:"to_lng"`
	DepartTime   string  `json:"depart_time"`
	ArriveTime   string  `json:"arrive_time"`
	Price        float64 `json:"price"`
	Notes        string  `json:"notes"`
}

func (h *MapHandler) GetMap(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	db := database.GetDB()

	var hotels []models.Hotel
	if err := db.Where("user_id = ?", userID).Find(&hotels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hotels"})
		return
	}

	var flights []models.Flight
	if err := db.Where("user_id = ?", userID).Find(&flights).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch flights"})
		return
	}

	hotelResp := make([]mapHotel, 0, len(hotels))
	for _, h := range hotels {
		hotelResp = append(hotelResp, mapHotel{
			ID:        h.ID,
			HotelName: h.HotelName,
			City:      h.City,
			Country:   h.Country,
			Latitude:  h.Latitude,
			Longitude: h.Longitude,
		})
	}

	flightResp := make([]mapFlight, 0, len(flights))
	for _, f := range flights {
		fromLat, fromLng, fromOK := lookupAirport(strings.ToUpper(f.FromAirport))
		toLat, toLng, toOK := lookupAirport(strings.ToUpper(f.ToAirport))
		// If not found, keep zero values; frontend can geocode as fallback.
		_ = fromOK
		_ = toOK

		flightResp = append(flightResp, mapFlight{
			ID:           f.ID,
			Airline:      f.Airline,
			FlightNumber: f.FlightNumber,
			FromAirport:  strings.ToUpper(f.FromAirport),
			ToAirport:    strings.ToUpper(f.ToAirport),
			FromLat:      fromLat,
			FromLng:      fromLng,
			ToLat:        toLat,
			ToLng:        toLng,
			DepartTime:   f.DepartTime.Format(time.RFC3339),
			ArriveTime:   f.ArriveTime.Format(time.RFC3339),
			Price:        f.Price,
			Notes:        f.Notes,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"hotels":  hotelResp,
		"flights": flightResp,
	})
}

// Minimal airport coordinate lookup for common samples; frontend can fall back if zero.
func lookupAirport(iata string) (float64, float64, bool) {
	coords := map[string][2]float64{
		"JFK": {40.6413, -73.7781},
		"LAX": {33.9416, -118.4085},
		"SFO": {37.6213, -122.3790},
		"ORD": {41.9742, -87.9073},
		"MIA": {25.7959, -80.2871},
		"BOS": {42.3656, -71.0096},
		"MCO": {28.4312, -81.3081},
		"ATL": {33.6407, -84.4277},
		"SEA": {47.4502, -122.3088},
		"DFW": {32.8998, -97.0403},
	}

	if v, ok := coords[iata]; ok {
		return v[0], v[1], true
	}
	return 0, 0, false
}
