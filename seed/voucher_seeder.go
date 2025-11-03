package seed

import (
	"log"
	"time"

	"VSA_GOGIN_BE/models"

	"gorm.io/gorm"
)

func SeedVouchers(db *gorm.DB) {
	// Define the default vouchers
	defaultVouchers := []models.Voucher{
		{
			CrewName:        "Sinta",
			CrewID:          "S001",
			FlightNumber:    "ID001",
			FlightDate:      time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
			AircraftType:    "ATR 72",
			AircraftTypeKey: "atr_72",
			Seat1:           "1F",
			Seat2:           "6A",
			Seat3:           "2D",
		},
		{
			CrewName:        "Dwi",
			CrewID:          "D001",
			FlightNumber:    "ID001",
			FlightDate:      time.Date(2025, 12, 2, 0, 0, 0, 0, time.UTC),
			AircraftType:    "Airbus A320",
			AircraftTypeKey: "airbus_a320",
			Seat1:           "8F",
			Seat2:           "3C",
			Seat3:           "17B",
		},
		// Add more vouchers as needed
	}

	for _, voucher := range defaultVouchers {
		var existing models.Voucher
		err := db.Where("flight_number = ? AND flight_date = ? AND aircraft_type_key = ? AND crew_id = ?", voucher.FlightNumber, voucher.FlightDate, voucher.AircraftTypeKey, voucher.CrewID).First(&existing).Error

		// Only create if not exists
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&voucher).Error; err != nil {
				log.Printf("Failed to seed voucher: %v", err)
			} else {
				log.Printf("Seeded voucher: %s on %s with aircraft type %s", voucher.FlightNumber, voucher.FlightDate, voucher.AircraftType)
			}
		}
	}
}
