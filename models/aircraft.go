package models

import (
	"strings"

	"gorm.io/gorm"
)

type Aircraft struct {
	ID              uint   `json:"id" gorm:"primaryKey"`
	AircraftType    string `json:"aircraft_type" gorm:"unique;not null"`
	AircraftTypeKey string `json:"aircraft_type_key" gorm:"unique;not null"`
	NumRows         int    `json:"num_rows"`
	SeatsPerRow     string `json:"seats_per_row"` // Comma-separated seat letters, e.g. "A,C,D,F"
}

// BeforeSave hook — automatically set AircraftTypeKey before saving
func (a *Aircraft) BeforeSave(tx *gorm.DB) (err error) {
	key := strings.ToLower(a.AircraftType)
	key = strings.ReplaceAll(key, " ", "_")
	a.AircraftTypeKey = key
	return
}
