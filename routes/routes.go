package routes

import (
	"github.com/gin-gonic/gin"
	"VSA_GOGIN_BE/controllers"
)

func SetupVoucherRoutes(router *gin.Engine, controller *controllers.VoucherController) {
	vouchers := router.Group("/api/vouchers")
	{
		vouchers.POST("/", controller.CreateVoucher)
		vouchers.GET("/", controller.ListVouchers)
		vouchers.GET("/:id", controller.GetVoucher)
		vouchers.PUT("/:id", controller.UpdateVoucher)
		vouchers.DELETE("/:id", controller.DeleteVoucher)
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

func SetupFlightRoutes(router *gin.Engine, controller *controllers.FlightController) {
	flights := router.Group("/api/flights")
	{
		flights.POST("/", controller.CreateFlight)
		flights.GET("/", controller.ListFlights)
		flights.GET("/:id", controller.GetFlight)
		flights.PUT("/:id", controller.UpdateFlight)
		flights.DELETE("/:id", controller.DeleteFlight)
	}
}