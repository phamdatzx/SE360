package router

import (
	"user-service/controller"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, c controller.UserController) {
    user := rg.Group("/users")
    {
		user.POST("/register", c.Register)
    }
}
