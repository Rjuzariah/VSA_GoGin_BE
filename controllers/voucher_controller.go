package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"VSA_GOGIN_BE/models"
)

type VoucherController struct {
	DB *gorm.DB
}

func NewVoucherController(db *gorm.DB) *VoucherController {
	return &VoucherController{DB: db}
}

// CreateVoucher godoc
// @Summary Create a new voucher
// @Description Create a new voucher with crew information and seat assignments
// @Tags vouchers
// @Accept json
// @Produce json
// @Param voucher body models.Voucher true "Voucher information"
// @Success 201 {object} models.Voucher "Successfully created voucher"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Server error"
// @Router /vouchers [post]
func (c *VoucherController) CreateVoucher(ctx *gin.Context) {
	var voucher models.Voucher
	if err := ctx.ShouldBindJSON(&voucher); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.DB.Create(&voucher).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, voucher)
}

// GetVoucher godoc
// @Summary Get a voucher by ID
// @Description Get detailed voucher information including crew details and seat assignments
// @Tags vouchers
// @Produce json
// @Param id path int true "Voucher ID"
// @Success 200 {object} models.Voucher "Successfully retrieved voucher"
// @Failure 404 {object} map[string]string "Voucher not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /vouchers/{id} [get]
func (c *VoucherController) GetVoucher(ctx *gin.Context) {
	var voucher models.Voucher
	if err := c.DB.First(&voucher, ctx.Param("id")).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Voucher not found"})
		return
	}

	ctx.JSON(http.StatusOK, voucher)
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

// UpdateVoucher godoc
// @Summary Update a voucher
// @Description Update voucher details including crew information and seat assignments
// @Tags vouchers
// @Accept json
// @Produce json
// @Param id path int true "Voucher ID"
// @Param voucher body models.Voucher true "Updated voucher information"
// @Success 200 {object} models.Voucher "Successfully updated voucher"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "Voucher not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /vouchers/{id} [put]
func (c *VoucherController) UpdateVoucher(ctx *gin.Context) {
	var voucher models.Voucher
	if err := c.DB.First(&voucher, ctx.Param("id")).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Voucher not found"})
		return
	}

	if err := ctx.ShouldBindJSON(&voucher); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.DB.Save(&voucher).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, voucher)
}

// DeleteVoucher godoc
// @Summary Delete a voucher
// @Description Delete a voucher and its associated data
// @Tags vouchers
// @Param id path int true "Voucher ID"
// @Success 204 "Successfully deleted"
// @Failure 500 {object} map[string]string "Server error"
// @Router /vouchers/{id} [delete]
func (c *VoucherController) DeleteVoucher(ctx *gin.Context) {
	if err := c.DB.Delete(&models.Voucher{}, ctx.Param("id")).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
