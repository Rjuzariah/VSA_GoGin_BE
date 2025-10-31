// @title Voucher Seat Assignment API
// @version 1.0
// @description This is a service for managing aircraft, flight, and voucher assignments
// @host localhost:8081
// @BasePath /api
package main

import (
	"VSA_GOGIN_BE/config"
	"VSA_GOGIN_BE/controllers"
	"VSA_GOGIN_BE/models"
	"VSA_GOGIN_BE/routes"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	// Swagger
	_ "VSA_GOGIN_BE/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Create a Gin router with default middleware
	router := gin.Default()

	redisClient := config.InitRedis()
	fmt.Println("✅ Connected to Redis")

	// Example: store a key
	redisClient.Set(config.Ctx, "test", "hello", 0)

	// Enable CORS for frontend development
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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
