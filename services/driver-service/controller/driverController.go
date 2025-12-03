package controller

import (
	"driver-service/dto"
	apperror "driver-service/error"
	"driver-service/service"
	"driver-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type DriverController struct {
	service service.DriverService
}

func NewDriverController(service service.DriverService) *DriverController {
	return &DriverController{service: service}
}

// UpdateLocation handles PUT /drivers/:id/location
func (c *DriverController) UpdateLocation(ctx *gin.Context) {
	driverID := ctx.Param("id")

	var request dto.UpdateLocationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate struct
	if err := validate.Struct(request); err != nil {
		var errors string
		for _, err := range err.(validator.ValidationErrors) {
			errors += err.Field() + " is invalid: " + err.Tag() + ", "
		}
		_ = ctx.Error(apperror.NewAppError(400, errors))
		ctx.Abort()
		return
	}

	result, err := c.service.UpdateDriverLocation(driverID, request)
	if err != nil {
		_ = ctx.Error(apperror.NewAppErrorWithErr(404, "Driver not found or update failed", err))
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Location updated successfully", result)
}

// SearchDrivers handles GET /drivers/search
func (c *DriverController) SearchDrivers(ctx *gin.Context) {
	var request dto.SearchDriversRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate struct
	if err := validate.Struct(request); err != nil {
		var errors string
		for _, err := range err.(validator.ValidationErrors) {
			errors += err.Field() + " is invalid: " + err.Tag() + ", "
		}
		_ = ctx.Error(apperror.NewAppError(400, errors))
		ctx.Abort()
		return
	}

	results, err := c.service.SearchDrivers(request)
	if err != nil {
		_ = ctx.Error(apperror.NewAppErrorWithErr(500, "Search failed", err))
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Search completed successfully", results)
}
