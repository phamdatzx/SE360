package main

import (
	"user-service/config"
	"user-service/controller"
	"user-service/initializer"
	"user-service/model"
	"user-service/repository"
	"user-service/router"
	"user-service/service"
)

func main() {
	initializer.LoadEnv()

	// Kết nối DB
	config.ConnectDatabase()

	// Auto migrate
	config.DB.AutoMigrate(&model.User{})

	//wiring dependencies
	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	// Setup routes
	r := router.SetupRouter(&router.AppRouter{
		UserController: userController,
	})

	r.Run(":8080") // chạy server ở port 8080
}
