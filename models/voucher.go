package models

import "time"

type Voucher struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	CrewName        string    `json:"crew_name"`
	CrewID          string    `json:"crew_id"`
	FlightNumber    string    `json:"flight_number"`
	FlightDate      time.Time `json:"flight_date"`
	AircraftType    string    `json:"aircraft_type"`
	AircraftTypeKey string    `json:"aircraft_type_key"`
	Seat1           string    `json:"seat1"`
	Seat2           string    `json:"seat2"`
	Seat3           string    `json:"seat3"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
}
