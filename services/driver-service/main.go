package main

import (
	"driver-service/config"
	"driver-service/controller"
	"driver-service/model"
	"driver-service/repository"
	"driver-service/router"
	"driver-service/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("🚗 Driver Service Starting...")

	// Connect to database
	config.ConnectDatabase()

	// Auto migrate
	config.DB.AutoMigrate(&model.Driver{})

	// Initialize repository
	driverRepo := repository.NewDriverRepository(config.DB)

	// Seed mock data
	if err := driverRepo.SeedMockData(); err != nil {
		fmt.Printf("⚠️  Warning: Failed to seed mock data: %v\n", err)
	}

	// Wire dependencies
	driverService := service.NewDriverService(driverRepo)
	driverController := controller.NewDriverController(driverService)

	// Setup Gin router
	r := gin.Default()
	//r.Use(cors.Default())

	// Setup routes
	router.SetupRouter(r, &router.AppRouter{
		DriverController: driverController,
	})

	// Run server
	fmt.Println("🚀 Driver Service running on port 8082")
	r.Run(":8082")
}
