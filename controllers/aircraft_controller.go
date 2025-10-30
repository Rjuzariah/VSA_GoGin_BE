package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"VSA_GOGIN_BE/models"
)

type AircraftController struct {
	DB *gorm.DB
}

func NewAircraftController(db *gorm.DB) *AircraftController {
	return &AircraftController{DB: db}
}

// CreateAircraft godoc
// @Summary Create a new aircraft
// @Description Create a new aircraft with the provided details
// @Tags aircraft
// @Accept json
// @Produce json
// @Success 201 {object} models.Aircraft
// @Router /aircraft [post]
func (c *AircraftController) CreateAircraft(ctx *gin.Context) {
	var aircraft models.Aircraft
	if err := ctx.ShouldBindJSON(&aircraft); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validation: NumRows must be > 0, SeatsPerRow must not be empty
	if aircraft.NumRows <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "num_rows must be greater than 0"})
		return
	}
	if aircraft.SeatsPerRow == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "seats_per_row must not be empty"})
		return
	}

	if err := c.DB.Create(&aircraft).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, aircraft)
}

// GetAircraft godoc
// @Summary Get an aircraft by ID
// @Description Get aircraft details by ID
// @Tags aircraft
// @Produce json
// @Param id path int true "Aircraft ID"
// @Success 200 {object} models.Aircraft
// @Router /aircraft/{id} [get]
func (c *AircraftController) GetAircraft(ctx *gin.Context) {
	var aircraft models.Aircraft
	if err := c.DB.First(&aircraft, ctx.Param("id")).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Aircraft not found"})
		return
	}

	ctx.JSON(http.StatusOK, aircraft)
}

// ListAircraft godoc
// @Summary List all aircraft
// @Description Get all aircraft
// @Tags aircraft
// @Produce json
// @Success 200 {array} models.Aircraft
// @Router /aircraft [get]
func (c *AircraftController) ListAircraft(ctx *gin.Context) {
	var aircraft []models.Aircraft
	if err := c.DB.Find(&aircraft).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, aircraft)
}

// UpdateAircraft godoc
// @Summary Update an aircraft
// @Description Update aircraft details by ID
// @Tags aircraft
// @Accept json
// @Produce json
// @Param id path int true "Aircraft ID"
// @Success 200 {object} models.Aircraft
// @Router /aircraft/{id} [put]
func (c *AircraftController) UpdateAircraft(ctx *gin.Context) {
	var aircraft models.Aircraft
	if err := c.DB.First(&aircraft, ctx.Param("id")).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Aircraft not found"})
		return
	}

	if err := ctx.ShouldBindJSON(&aircraft); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validation: NumRows must be > 0, SeatsPerRow must not be empty
	if aircraft.NumRows <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "num_rows must be greater than 0"})
		return
	}
	if aircraft.SeatsPerRow == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "seats_per_row must not be empty"})
		return
	}

	if err := c.DB.Save(&aircraft).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, aircraft)
}

// DeleteAircraft godoc
// @Summary Delete an aircraft
// @Description Delete aircraft by ID
// @Tags aircraft
// @Param id path int true "Aircraft ID"
// @Success 204 "No Content"
// @Router /aircraft/{id} [delete]
func (c *AircraftController) DeleteAircraft(ctx *gin.Context) {
	if err := c.DB.Delete(&models.Aircraft{}, ctx.Param("id")).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}