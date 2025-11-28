package controller

import (
	"net/http"
	"strings"
	"user-service/dto"
	appError "user-service/error"
	"user-service/service"
	"user-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type UserController struct {
	service service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) Register(ctx *gin.Context) {
	var request dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ validate struct
	if err := validate.Struct(request); err != nil {
		// lấy chi tiết lỗi
		var errors string
		for _, err := range err.(validator.ValidationErrors) {
			errors += err.Field() + " is invalid: " + err.Tag() + ", "
		}

		_ = ctx.Error(appError.NewAppError(400, errors))
		ctx.Abort()
		return
	}

	resultUser, err := c.service.RegisterUser(request)
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

// auth-forward verify controller
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
	ctx.Header("X-Username", claims.Username)
	ctx.Header("X-User-Role", claims.Role)

	utils.SuccessResponse(ctx, 200, "Verify successfully", nil)
}

func (c *UserController) ActivateAccount(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		ctx.Error(appError.NewAppError(401, "Token is required"))
		ctx.Abort()
		return
	}

	if err := c.service.ActivateAccount(token); err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Activate account successfully", nil)
}

func (c *UserController) SendResetPasswordRequest(ctx *gin.Context) {
	var request dto.SendResetPasswordRequestDto
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := c.service.SendResetPasswordRequest(request); err != nil {
		_ = ctx.Error(appError.NewAppError(400, err.Error()))
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "SendResetPasswordRequest successfully", nil)
}

func (c *UserController) ResetPassword(ctx *gin.Context) {
	//get token from param
	token := ctx.Query("token")
	if token == "" {
		ctx.Error(appError.NewAppError(401, "Token is required"))
		ctx.Abort()
		return
	}

	//get body
	var request dto.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.ResetPassword(token, request); err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Reset account password successfully", nil)
}

func (c *UserController) TestPrivate(ctx *gin.Context) {
	utils.SuccessResponse(ctx, 200, "Reached", nil)
}
