package models

import "time"

type Aircraft struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	AircraftType string    `json:"aircraft_type"`
	NumRows      int       `json:"num_rows"`
	SeatsPerRow  string    `json:"seats_per_row"` // Comma-separated seat letters, e.g. "A,C,D,F"
	CreatedAt    time.Time `json:"created_at"`
}
