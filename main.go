// @title Voucher Seat Assignment API
// @version 1.0
// @description This is a service for managing aircraft, flight, and voucher assignments
// @host localhost:8081
// @BasePath /api
package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"VSA_GOGIN_BE/controllers"
	"VSA_GOGIN_BE/models"
	"VSA_GOGIN_BE/routes"

	// Swagger
	_ "VSA_GOGIN_BE/docs"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
)

func main() {
	// Create a Gin router with default middleware
	router := gin.Default()

	// Database connection using SQLite
	db, err := gorm.Open(sqlite.Open("vsa.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.Aircraft{}, &models.Flight{}, &models.Voucher{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize controllers
	aircraftController := controllers.NewAircraftController(db)
	flightController := controllers.NewFlightController(db)
	voucherController := controllers.NewVoucherController(db)

	// Setup routes
	routes.SetupAircraftRoutes(router, aircraftController)
	routes.SetupFlightRoutes(router, flightController)
	routes.SetupVoucherRoutes(router, voucherController)

	// Swagger UI endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Basic health check endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Start the server
    log.Println("Server starting on :8081")
    if err := router.Run(":8081"); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}