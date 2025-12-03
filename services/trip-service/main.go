package main

import (
	"fmt"
	"trip-service/config"
	"trip-service/controller"
	"trip-service/model"
	"trip-service/repository"
	"trip-service/router"
	"trip-service/service"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("🚕 Trip Service Starting...")

	// Connect to database
	config.ConnectDatabase()

	// Auto migrate
	config.DB.AutoMigrate(&model.Trip{})

	// Wire dependencies
	tripRepo := repository.NewTripRepository(config.DB)
	tripService := service.NewTripService(tripRepo)
	tripController := controller.NewTripController(tripService)

	// Setup Gin router
	r := gin.Default()
	//r.Use(cors.Default())

	// Setup routes
	router.SetupRouter(r, &router.AppRouter{
		TripController: tripController,
	})

	// Run server
	fmt.Println("🚀 Trip Service running on port 8083")
	r.Run(":8083")
}
