package controller

import (
	"net/http"
	"strings"
	"user-service/dto"
	appError "user-service/error"
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

func (c *UserController) Verify(ctx *gin.Context) {
	//get auth header
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.Error(appError.NewAppError(401, "Authorization header not provided"))
		ctx.Abort()
		return
	}

	//remove "Bearer" prefix
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		ctx.Error(appError.NewAppError(401, "Invalid Authorization header format"))
		ctx.Abort()
		return
	}
	token := parts[1]

	claims, err := utils.VerifyToken(token)
	if err != nil {
		ctx.Error(appError.NewAppErrorWithErr(401, "Invalid token", err))
		ctx.Abort()
		return
	}

	// Trả header để Traefik forward sang backend
	ctx.Header("X-User-Id", claims.UserID)
	ctx.Header("X-User-Role", claims.Role)

	utils.SuccessResponse(ctx, 200, "Verify successfully", nil)
}
