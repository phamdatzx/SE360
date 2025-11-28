package router

import (
	"product-service/controller"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(rg *gin.RouterGroup, c controller.ProductController) {
	product := rg.Group("")
	{
		product.POST("/", c.CreateProduct)
		product.GET("/:id", c.GetProductByID)
		product.GET("/", c.GetAllProducts)
		product.PUT("/:id", c.UpdateProduct)
		product.DELETE("/:id", c.DeleteProduct)
	}
}
