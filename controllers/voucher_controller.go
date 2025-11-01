package controllers

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"VSA_GOGIN_BE/models"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type VoucherController struct {
	DB  *gorm.DB
	RDB *redis.Client // Add Redis client
}

func NewVoucherController(db *gorm.DB, rdb *redis.Client) *VoucherController {
	return &VoucherController{
		DB:  db,
		RDB: rdb,
	}
}

// ListVouchers godoc
// @Summary List all vouchers
// @Description Get a list of all vouchers with their associated flight information
// @Tags vouchers
// @Produce json
// @Success 200 {array} models.Voucher "Successfully retrieved vouchers list"
// @Failure 500 {object} map[string]string "Server error"
// @Router /vouchers [get]
func (c *VoucherController) ListVouchers(ctx *gin.Context) {
	var vouchers []models.Voucher
	if err := c.DB.Find(&vouchers).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, vouchers)
}

func (c *VoucherController) CheckVoucherSeat(ctx *gin.Context) {
	var voucher models.Voucher

	if err := ctx.ShouldBindJSON(&voucher); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if voucher has already been generated
	var vouchers []models.Voucher
	if err := c.DB.Where("flight_number = ?", voucher.FlightNumber).Where("flight_date = ?", voucher.FlightDate).Where("crew_id = ?", voucher.CrewID).Find(&vouchers).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	if len(vouchers) > 0 {
		ctx.JSON(http.StatusOK, gin.H{"exists": true})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"exists": false})
	}

}

// @Summary Generate a random voucher seat
// @Description Generate a random seat for a voucher based on the flight ID
// @Tags vouchers
// @Produce json
// @Param flight_id path string true "Flight ID"
// @Success 200 {object} string "Random seat"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Server error"
// @Router /vouchers/generate [post]
func (c *VoucherController) GenerateVoucherSeat(ctx *gin.Context) {
	var voucher models.Voucher

	if err := ctx.ShouldBindJSON(&voucher); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if voucher has already been generated
	var vouchers []models.Voucher
	if err := c.DB.Where("flight_number = ?", voucher.FlightNumber).Where("flight_date = ?", voucher.FlightDate).Where("crew_id = ?", voucher.CrewID).Find(&vouchers).Error; err != nil {
		return
	}
	if len(vouchers) > 0 {
		ctx.JSON(http.StatusConflict, gin.H{"message": "Voucher already generated for this crew member on this flight and date"})
		return
	}

	cacheKey := "voucher_seat_cache:" + voucher.FlightNumber + ":" + voucher.FlightDate.String()

	// 1️⃣ Try to get cache from Redis
	var allSeats []string
	cacheSeat, err := c.RDB.Get(context.Background(), cacheKey).Result()
	json.Unmarshal([]byte(cacheSeat), &allSeats)

	if err != nil {
		fmt.Println("Cache not found")
	}

	fmt.Print("Hello, world")

	if allSeats == nil {
		//GENERATE all seats
		fmt.Println("Generating all available seats for flight:", voucher.FlightNumber)
		//get row and seat info from aircraft associated with flightID
		var flight models.Flight
		if err := c.DB.Preload("Aircraft").First(&flight, "flight_number = ?", voucher.FlightNumber).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
			return
		}

		//get seat1,seat2,seat3 from vouchers associated with flightID to skip
		var vouchers []models.Voucher
		if err := c.DB.Where("flight_number = ?", voucher.FlightNumber).Where("flight_date = ?", voucher.FlightDate).Find(&vouchers).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//map voucher seats to skip list
		skipList := map[string]bool{}
		for _, voucher := range vouchers {
			if voucher.Seat1 != "" {
				skipList[voucher.Seat1] = true
			}
			if voucher.Seat2 != "" {
				skipList[voucher.Seat2] = true
			}
			if voucher.Seat3 != "" {
				skipList[voucher.Seat3] = true
			}
		}

		rows := flight.Aircraft.NumRows
		sections := flight.Aircraft.SeatsPerRow // e.g. "ABCDEF"
		fmt.Println("Rows:", rows)
		fmt.Println("Sections:", sections)
		//generate all seats excluding skipList
		for i := 1; i <= rows; i++ {
			for _, section := range sections {
				seat := fmt.Sprintf("%d%c", i, section)
				if !skipList[seat] {
					allSeats = append(allSeats, seat)
				}
			}
		}
		//end generate all seats

		//store allSeat to Redis cache
		seatsJSON, _ := json.Marshal(allSeats)
		c.RDB.Set(context.Background(), cacheKey, seatsJSON, time.Hour*1)
	}

	// Pick a random available seat
	// shuffle allSeats
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(allSeats), func(i, j int) {
		allSeats[i], allSeats[j] = allSeats[j], allSeats[i]
	})

	// get first 3 random seats (ensure there are at least 3 available)
	numSeats := 3
	if len(allSeats) < numSeats {
		numSeats = len(allSeats)
	}
	randomSeats := allSeats[:numSeats]

	// create voucherseat data
	// TO DO: Save the assigned random seats to voucher and update Redis cache accordingly
	// Assign selected seats to voucher
	voucher.Seat1 = randomSeats[0]
	voucher.Seat2 = randomSeats[1]
	voucher.Seat3 = randomSeats[2]
	if err := c.DB.Save(&voucher).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(randomSeats) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No seats found for this flight"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":     true,
		"random_seat": randomSeats,
	})
}
