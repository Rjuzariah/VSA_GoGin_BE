package seed

import (
	"log"

	"VSA_GOGIN_BE/models"

	"gorm.io/gorm"
)

func SeedAircrafts(db *gorm.DB) {
	// Define the default aircrafts
	defaultAircrafts := []models.Aircraft{
		{AircraftTypeKey: "atr_72", AircraftType: "ATR 72", NumRows: 18, SeatsPerRow: "ABCDF"},
		{AircraftTypeKey: "airbus_a320", AircraftType: "Airbus A320", NumRows: 32, SeatsPerRow: "ABCDEF"},
		{AircraftTypeKey: "boeing_737", AircraftType: "Boeing 737", NumRows: 32, SeatsPerRow: "ABCDEF"},
	}

	for _, aircraft := range defaultAircrafts {
		var existing models.Aircraft
		err := db.Where("aircraft_type = ?", aircraft.AircraftType).First(&existing).Error

		// Only create if not exists
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&aircraft).Error; err != nil {
				log.Printf("Failed to seed aircraft %s: %v", aircraft.AircraftType, err)
			} else {
				log.Printf("Seeded aircraft: %s", aircraft.AircraftType)
			}
		}
	}
}
