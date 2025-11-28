package main

import (
	"fmt"
	"user-service/config"
	"user-service/controller"
	"user-service/model"
	"user-service/repository"
	"user-service/router"
	"user-service/service"

	"github.com/gin-gonic/gin"
)

func main() {
	//

	fmt.Println("Hello World")
	// Kết nối DB
	config.ConnectDatabase()

	// Auto migrate
	config.DB.AutoMigrate(&model.User{})

	//wiring dependencies
	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	r := gin.Default()
	//r.Use(cors.Default())

	// Setup routes
	router.SetupRouter(r, &router.AppRouter{
		UserController: userController,
	})

	r.Run(":8080") // chạy server ở port 8080
}
