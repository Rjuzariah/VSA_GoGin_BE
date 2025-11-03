package routes

import (
	"VSA_GOGIN_BE/controllers"

	"github.com/gin-gonic/gin"
)

func SetupVoucherRoutes(router *gin.Engine, controller *controllers.VoucherController) {
	vouchers := router.Group("/api/vouchers")
	{
		vouchers.GET("/", controller.ListVouchers)
		vouchers.POST("/generate", controller.GenerateVoucherSeat)
		vouchers.POST("/check", controller.CheckVoucherSeat)
	}
}

func SetupAircraftRoutes(router *gin.Engine, controller *controllers.AircraftController) {
	aircraft := router.Group("/api/aircraft")
	{
		aircraft.POST("/", controller.CreateAircraft)
		aircraft.GET("/", controller.ListAircraft)
		aircraft.GET("/:id", controller.GetAircraft)
		aircraft.PUT("/:id", controller.UpdateAircraft)
		aircraft.DELETE("/:id", controller.DeleteAircraft)
	}
}
