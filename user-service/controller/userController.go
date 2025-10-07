package controller

import (
	"net/http"
	"user-service/model"
	"user-service/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) Register(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resultUser, err := c.service.RegisterUser(user)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, resultUser)
}
