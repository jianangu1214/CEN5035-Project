package models

import (
	"time"

	"gorm.io/gorm"
)

type Hotel struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	HotelName   string         `gorm:"not null" json:"hotel_name"`
	City        string         `gorm:"not null" json:"city"`
	Country     string         `gorm:"not null" json:"country"`
	CheckIn     time.Time      `gorm:"not null" json:"check_in"`
	CheckOut    time.Time      `gorm:"not null" json:"check_out"`
	Price       float64        `gorm:"not null" json:"price"`
	Notes       string         `json:"notes"`
	Latitude    float64        `json:"latitude"`
	Longitude   float64        `json:"longitude"`
}

type HotelResponse struct {
	ID        uint    `json:"id"`
	HotelName string  `json:"hotel_name"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
	CheckIn   string  `json:"check_in"`
	CheckOut  string  `json:"check_out"`
	Price     float64 `json:"price"`
	Notes     string  `json:"notes"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CreateHotelRequest struct {
	HotelName string  `json:"hotel_name" binding:"required"`
	City      string  `json:"city" binding:"required"`
	Country   string  `json:"country" binding:"required"`
	CheckIn   string  `json:"check_in" binding:"required"`
	CheckOut  string  `json:"check_out" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
	Notes     string  `json:"notes"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type UpdateHotelRequest struct {
	HotelName string  `json:"hotel_name"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
	CheckIn   string  `json:"check_in"`
	CheckOut  string  `json:"check_out"`
	Price     float64 `json:"price"`
	Notes     string  `json:"notes"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (Hotel) TableName() string {
	return "hotels"
}

func (h *Hotel) ToResponse() HotelResponse {
	return HotelResponse{
		ID:        h.ID,
		HotelName: h.HotelName,
		City:      h.City,
		Country:   h.Country,
		CheckIn:   h.CheckIn.Format("2006-01-02"),
		CheckOut:  h.CheckOut.Format("2006-01-02"),
		Price:     h.Price,
		Notes:     h.Notes,
		Latitude:  h.Latitude,
		Longitude: h.Longitude,
	}
}
