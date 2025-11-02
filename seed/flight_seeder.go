package seed

import (
	"log"
	"time"

	"VSA_GOGIN_BE/models"

	"gorm.io/gorm"
)

func SeedFlights(db *gorm.DB) {
	// Define the default flights
	flight1Date, _ := time.Parse(time.RFC3339, "2025-12-01T00:00:00Z")
	flight2Date, _ := time.Parse(time.RFC3339, "2025-12-02T00:00:00Z")
	flight3Date, _ := time.Parse(time.RFC3339, "2025-12-03T00:00:00Z")

	defaultFlights := []models.Flight{
		{FlightNumber: "ID001", FlightDate: flight1Date, AircraftID: 1},
		{FlightNumber: "ID002", FlightDate: flight2Date, AircraftID: 2},
		{FlightNumber: "ID003", FlightDate: flight3Date, AircraftID: 3},
	}

	for _, flight := range defaultFlights {
		var existing models.Flight
		err := db.Where("flight_number = ? AND DATE(flight_date) = DATE(?)", flight.FlightNumber, flight.FlightDate).
			First(&existing).Error

		// Only create if not exists
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&flight).Error; err != nil {
				log.Printf("Failed to seed flight %s: %v", flight.FlightNumber, err)
			} else {
				log.Printf("Seeded flight: %s (%s)", flight.FlightNumber, flight.FlightDate.Format("2006-01-02"))
			}
		}
	}
}
