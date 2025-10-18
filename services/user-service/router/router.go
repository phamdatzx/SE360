package router

import (
	"user-service/controller"
	"user-service/middleware"

	"github.com/gin-gonic/gin"
)

// AppRouter holds all controllers for dependency injection
type AppRouter struct {
	UserController *controller.UserController
}

// SetupRouter builds the main Gin router and registers all module routes
func SetupRouter(appRouter *AppRouter) *gin.Engine {
	engine := gin.Default()

	// Global middlewares
	engine.Use(
		middleware.ErrorHandler(),
	)

	api := engine.Group("/api/user")
	{
		// Each controller registers its own routes
		RegisterUserRoutes(api, *appRouter.UserController)
	}

	return engine
}
