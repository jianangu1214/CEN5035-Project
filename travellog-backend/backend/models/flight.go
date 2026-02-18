package models

import "time"

type Flight struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	Airline      string    `json:"airline"`
	FlightNumber string    `json:"flight_number"`
	FromAirport  string    `gorm:"size:3" json:"from_airport"`
	ToAirport    string    `gorm:"size:3" json:"to_airport"`
	DepartTime   time.Time `json:"depart_time"`
	ArriveTime   time.Time `json:"arrive_time"`
	Price        float64   `json:"price"`
	Notes        string    `gorm:"type:text" json:"notes"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
