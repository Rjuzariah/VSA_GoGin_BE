package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"VSA_GOGIN_BE/models"
)

type FlightController struct {
	DB *gorm.DB
}

func NewFlightController(db *gorm.DB) *FlightController {
	return &FlightController{DB: db}
}

// CreateFlight godoc
// @Summary Create a new flight
// @Description Create a new flight with the provided details
// @Tags flights
// @Accept json
// @Produce json
// @Success 201 {object} models.Flight
// @Router /flights [post]
func (c *FlightController) CreateFlight(ctx *gin.Context) {
	var flight models.Flight
	if err := ctx.ShouldBindJSON(&flight); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify that the referenced aircraft exists
	var aircraft models.Aircraft
	if err := c.DB.First(&aircraft, flight.AircraftID).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Referenced aircraft not found"})
		return
	}

	if err := c.DB.Create(&flight).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, flight)
}

// GetFlight godoc
// @Summary Get a flight by ID
// @Description Get flight details by ID
// @Tags flights
// @Produce json
// @Param id path int true "Flight ID"
// @Success 200 {object} models.Flight
// @Router /flights/{id} [get]
func (c *FlightController) GetFlight(ctx *gin.Context) {
	var flight models.Flight
	if err := c.DB.Preload("Aircraft").First(&flight, ctx.Param("id")).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	ctx.JSON(http.StatusOK, flight)
}

// ListFlights godoc
// @Summary List all flights
// @Description Get all flights
// @Tags flights
// @Produce json
// @Success 200 {array} models.Flight
// @Router /flights [get]
func (c *FlightController) ListFlights(ctx *gin.Context) {
	var flights []models.Flight
	if err := c.DB.Preload("Aircraft").Find(&flights).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, flights)
}

// UpdateFlight godoc
// @Summary Update a flight
// @Description Update flight details by ID
// @Tags flights
// @Accept json
// @Produce json
// @Param id path int true "Flight ID"
// @Success 200 {object} models.Flight
// @Router /flights/{id} [put]
func (c *FlightController) UpdateFlight(ctx *gin.Context) {
	var flight models.Flight
	if err := c.DB.First(&flight, ctx.Param("id")).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	if err := ctx.ShouldBindJSON(&flight); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify that the referenced aircraft exists
	var aircraft models.Aircraft
	if err := c.DB.First(&aircraft, flight.AircraftID).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Referenced aircraft not found"})
		return
	}

	if err := c.DB.Save(&flight).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, flight)
}

// DeleteFlight godoc
// @Summary Delete a flight
// @Description Delete flight by ID
// @Tags flights
// @Param id path int true "Flight ID"
// @Success 204 "No Content"
// @Router /flights/{id} [delete]
func (c *FlightController) DeleteFlight(ctx *gin.Context) {
	if err := c.DB.Delete(&models.Flight{}, ctx.Param("id")).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}