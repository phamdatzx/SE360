package router

import (
	"user-service/controller"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, c controller.UserController) {
	user := rg.Group("")
	{
		user.POST("/public/register", c.Register)
		user.POST("/public/login", c.Login)
		user.GET("/public/verify", c.Verify)
		user.POST("/public/activate", c.ActivateAccount)
		user.POST("/public/send-reset-password", c.SendResetPasswordRequest)
		user.POST("/public/reset-password", c.ResetPassword)
		user.GET("/test-private", c.TestPrivate)
	}
}
