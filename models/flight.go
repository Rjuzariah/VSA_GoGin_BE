package models

import "time"

type Flight struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FlightNumber string   `json:"flight_number"`
	FlightDate   time.Time `json:"flight_date"`
	AircraftID   uint     `json:"aircraft_id"`
	CreatedAt    time.Time `json:"created_at"`
	Aircraft     Aircraft  `json:"aircraft" gorm:"foreignKey:AircraftID"`
}