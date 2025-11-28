package main

import (
	"fmt"
	"product-service/config"
	"product-service/controller"
	"product-service/repository"
	"product-service/router"
	"product-service/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello Product Service")
	// Kết nối DB
	config.ConnectDatabase()

	// Initialize S3
	config.InitS3()

	// MongoDB doesn't require AutoMigrate like GORM
	// Collections will be created automatically when first document is inserted

	//wiring dependencies
	productRepo := repository.NewProductRepository(config.DB)
	productService := service.NewProductService(productRepo)
	productController := controller.NewProductController(productService)

	r := gin.Default()
	r.Use(cors.Default())

	// Setup routes
	router.SetupRouter(r, &router.AppRouter{
		ProductController: productController,
	})

	r.Run(":8081") // chạy server ở port 8081 để tránh conflict với user-service
}
