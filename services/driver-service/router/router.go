package router

import (
	"driver-service/controller"
	"driver-service/middleware"

	"github.com/gin-gonic/gin"
)

type AppRouter struct {
	DriverController *controller.DriverController
}

func SetupRouter(r *gin.Engine, appRouter *AppRouter) {
	// Apply error handler middleware
	r.Use(middleware.ErrorHandler())

	// Health check endpoint
	r.GET("/api/driver/public/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "driver-service",
		})
	})

	// API group
	api := r.Group("/api/driver")
	{
		RegisterDriverRoutes(api, appRouter.DriverController)
	}
}

