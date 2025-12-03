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
func SetupRouter(engine *gin.Engine, appRouter *AppRouter) *gin.Engine {

	// Global middlewares
	engine.Use(
		middleware.ErrorHandler(),
	)

	// Health check endpoint
	engine.GET("/api/user/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "user-service",
		})
	})

	api := engine.Group("/api/user")
	{
		// Each controller registers its own routes
		RegisterUserRoutes(api, *appRouter.UserController)
	}

	//publicApi := engine.Group("/api/public")
	//{
	//	RegisterPublicUserRoutes(publicApi, *appRouter.UserController)
	//}

	return engine
}

