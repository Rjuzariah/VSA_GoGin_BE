package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"VSA_GOGIN_BE/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type VoucherService struct {
	DB  *gorm.DB
	RDB *redis.Client
}

// Check if voucher exists
func (s *VoucherService) CheckVoucherExists(voucher *models.Voucher) (bool, error) {
	var vouchers []models.Voucher
	if err := s.DB.
		Where("flight_number = ?", voucher.FlightNumber).
		Where("flight_date = ?", voucher.FlightDate).
		Where("crew_id = ?", voucher.CrewID).
		Where("aircraft_type_key = ?", voucher.AircraftTypeKey).
		Find(&vouchers).Error; err != nil {
		return false, errors.New("database error while checking voucher")
	}

	return len(vouchers) > 0, nil
}

// Generate voucher seats
func (s *VoucherService) GenerateVoucherSeats(voucher *models.Voucher) ([]string, error) {
	// 1️⃣ Check if already exists
	exists, err := s.CheckVoucherExists(voucher)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("Voucher already generated for crew id: %s on this flight and date", voucher.CrewID)
	}

	// 2️⃣ Try cache
	cacheKey := "voucher_seat_cache:" + voucher.FlightNumber + ":" + voucher.FlightDate.String() + ":" + voucher.AircraftType
	var allSeats []string

	cacheSeat, err := s.RDB.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(cacheSeat), &allSeats)
	}

	var aircraft models.Aircraft
	if err := s.DB.
		Where("aircraft_type_key = ?", voucher.AircraftTypeKey).
		First(&aircraft).Error; err != nil {
		return nil, errors.New("aircraft not found")
	}

	// 3️⃣ Generate if cache empty
	if len(allSeats) == 0 {
		var vouchers []models.Voucher
		if err := s.DB.
			Where("flight_number = ?", voucher.FlightNumber).
			Where("flight_date = ?", voucher.FlightDate).
			Find(&vouchers).Error; err != nil {
			return nil, errors.New("failed to load voucher data")
		}

		skipList := map[string]bool{}
		for _, v := range vouchers {
			for _, seat := range []string{v.Seat1, v.Seat2, v.Seat3} {
				if seat != "" {
					skipList[seat] = true
				}
			}
		}

		rows := aircraft.NumRows
		sections := aircraft.SeatsPerRow
		for i := 1; i <= rows; i++ {
			for _, section := range sections {
				seat := fmt.Sprintf("%d%c", i, section)
				if !skipList[seat] {
					allSeats = append(allSeats, seat)
				}
			}
		}

		// cache new seats
		seatsJSON, _ := json.Marshal(allSeats)
		s.RDB.Set(context.Background(), cacheKey, seatsJSON, time.Hour)
	}

	// 4️⃣ Randomly pick seats
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(allSeats), func(i, j int) {
		allSeats[i], allSeats[j] = allSeats[j], allSeats[i]
	})
	numSeats := 3
	if len(allSeats) < numSeats {
		numSeats = len(allSeats)
	}
	if numSeats == 0 {
		return nil, errors.New("no seats available for this flight")
	}

	selected := allSeats[:numSeats]
	voucher.AircraftType = aircraft.AircraftType
	voucher.Seat1 = selected[0]
	if len(selected) > 1 {
		voucher.Seat2 = selected[1]
	}
	if len(selected) > 2 {
		voucher.Seat3 = selected[2]
	}

	// 5️⃣ Save to DB
	if err := s.DB.Save(voucher).Error; err != nil {
		return nil, errors.New("failed to save voucher with assigned seats")
	}

	return selected, nil
}
