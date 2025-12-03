package router

import (
	"trip-service/controller"

	"github.com/gin-gonic/gin"
)

func RegisterTripRoutes(rg *gin.RouterGroup, c *controller.TripController) {
	trips := rg.Group("/public/trips")
	{
		trips.POST("", c.CreateTrip)
		trips.GET("/:id", c.GetTrip)
		trips.POST("/:id/cancel", c.CancelTrip)
		trips.PUT("/:id/status", c.UpdateTripStatus)
	}
}
