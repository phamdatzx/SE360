package controller

import (
	"net/http"
	"user-service/dto"
	"user-service/model"
	"user-service/service"
	"user-service/utils"

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

	utils.SuccessResponse(ctx, 200, "Register success", resultUser)
}

func (c *UserController) Login(ctx *gin.Context) {
	var request dto.LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := c.service.Login(request)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Login successfully", response)

}
