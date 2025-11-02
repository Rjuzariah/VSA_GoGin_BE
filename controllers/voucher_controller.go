package controllers

import (
	"VSA_GOGIN_BE/models"
	"VSA_GOGIN_BE/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type VoucherController struct {
	DB      *gorm.DB
	RDB     *redis.Client // Add Redis client
	Service *services.VoucherService
}

func NewVoucherController(db *gorm.DB, rdb *redis.Client) *VoucherController {
	service := &services.VoucherService{
		DB:  db,
		RDB: rdb,
	}

	return &VoucherController{
		DB:      db,
		RDB:     rdb,
		Service: service,
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

// Check voucher godoc
// @Summary Check voucher seat
// @Description Generate voucher seat for crew members based on the flight ID and flight date and return exist true/false
// @Tags vouchers
// @Produce json
// @Param request body map[string]interface{} true "Voucher generate request"
// @Success 200 {object} map[string]interface{} "List voucher seat"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Server Error"
// @Router /vouchers/check [post]
func (c *VoucherController) CheckVoucherSeat(ctx *gin.Context) {
	var voucher models.Voucher
	if err := ctx.ShouldBindJSON(&voucher); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := c.Service.CheckVoucherExists(&voucher)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"exists": exists})
}

// @Summary Generate voucher seats
// @Description Generate voucher seat for crew members based on the flight ID and flight date
// @Tags vouchers
// @Produce json
// @Param request body map[string]interface{} true "Voucher generate request"
// @Success 201 {object} map[string]interface{} "Example: {\"success\": true, \"seats\": [\"3B\", \"7C\", \"14D\"]}"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Server Error"
// @Router /vouchers/generate [post]
func (c *VoucherController) GenerateVoucherSeat(ctx *gin.Context) {
	var voucher models.Voucher
	if err := ctx.ShouldBindJSON(&voucher); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	seats, err := c.Service.GenerateVoucherSeats(&voucher)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"seats":   seats,
	})
}
