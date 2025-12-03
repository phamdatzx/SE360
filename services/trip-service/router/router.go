package router

import (
	"trip-service/controller"
	"trip-service/middleware"

	"github.com/gin-gonic/gin"
)

type AppRouter struct {
	TripController *controller.TripController
}

func SetupRouter(r *gin.Engine, appRouter *AppRouter) {
	// Apply error handler middleware
	r.Use(middleware.ErrorHandler())

	// Health check endpoint
	r.GET("/api/trip/public/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "trip-service",
		})
	})

	// API group
	api := r.Group("/api/trip")
	{
		RegisterTripRoutes(api, appRouter.TripController)
	}
}

