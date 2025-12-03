package router

import (
	"driver-service/controller"

	"github.com/gin-gonic/gin"
)

func RegisterDriverRoutes(rg *gin.RouterGroup, c *controller.DriverController) {
	drivers := rg.Group("/public/drivers")
	{
		drivers.PUT("/:id/location", c.UpdateLocation)
		drivers.GET("/search", c.SearchDrivers)
	}
}
