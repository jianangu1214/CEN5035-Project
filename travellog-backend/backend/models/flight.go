package models

import (
	"time"

	"gorm.io/gorm"
)

type Flight struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	Airline      string         `gorm:"not null" json:"airline"`
	FlightNumber string         `gorm:"not null" json:"flight_number"`
	FromAirport  string         `gorm:"size:3;not null" json:"from_airport"`
	ToAirport    string         `gorm:"size:3;not null" json:"to_airport"`
	DepartTime   time.Time      `gorm:"not null" json:"depart_time"`
	ArriveTime   time.Time      `gorm:"not null" json:"arrive_time"`
	Price        float64        `gorm:"not null" json:"price"`
	Notes        string         `gorm:"type:text" json:"notes"`
}

type FlightResponse struct {
	ID           uint    `json:"id"`
	Airline      string  `json:"airline"`
	FlightNumber string  `json:"flight_number"`
	FromAirport  string  `json:"from_airport"`
	ToAirport    string  `json:"to_airport"`
	DepartTime   string  `json:"depart_time"`
	ArriveTime   string  `json:"arrive_time"`
	Price        float64 `json:"price"`
	Notes        string  `json:"notes"`
}

type CreateFlightRequest struct {
	Airline      string  `json:"airline" binding:"required"`
	FlightNumber string  `json:"flight_number" binding:"required"`
	FromAirport  string  `json:"from_airport" binding:"required"`
	ToAirport    string  `json:"to_airport" binding:"required"`
	DepartTime   string  `json:"depart_time" binding:"required"`
	ArriveTime   string  `json:"arrive_time" binding:"required"`
	Price        float64 `json:"price" binding:"required"`
	Notes        string  `json:"notes"`
}

type UpdateFlightRequest struct {
	Airline      string  `json:"airline"`
	FlightNumber string  `json:"flight_number"`
	FromAirport  string  `json:"from_airport"`
	ToAirport    string  `json:"to_airport"`
	DepartTime   string  `json:"depart_time"`
	ArriveTime   string  `json:"arrive_time"`
	Price        float64 `json:"price"`
	Notes        string  `json:"notes"`
}

func (Flight) TableName() string {
	return "flights"
}

func (f *Flight) ToResponse() FlightResponse {
	return FlightResponse{
		ID:           f.ID,
		Airline:      f.Airline,
		FlightNumber: f.FlightNumber,
		FromAirport:  f.FromAirport,
		ToAirport:    f.ToAirport,
		DepartTime:   f.DepartTime.Format(time.RFC3339),
		ArriveTime:   f.ArriveTime.Format(time.RFC3339),
		Price:        f.Price,
		Notes:        f.Notes,
	}
}
