package controller

import (
	"net/http"
	"trip-service/dto"
	apperror "trip-service/error"
	"trip-service/service"
	"trip-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type TripController struct {
	service service.TripService
}

func NewTripController(service service.TripService) *TripController {
	return &TripController{service: service}
}

// CreateTrip handles POST /trips
func (c *TripController) CreateTrip(ctx *gin.Context) {
	var request dto.CreateTripRequest
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

	result, err := c.service.CreateTrip(request)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 201, "Trip created successfully", result)
}

// GetTrip handles GET /trips/:id
func (c *TripController) GetTrip(ctx *gin.Context) {
	tripID := ctx.Param("id")

	result, err := c.service.GetTrip(tripID)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Trip retrieved successfully", result)
}

// CancelTrip handles POST /trips/:id/cancel
func (c *TripController) CancelTrip(ctx *gin.Context) {
	tripID := ctx.Param("id")

	var request dto.CancelTripRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.service.CancelTrip(tripID, request)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Trip cancelled successfully", result)
}

// UpdateTripStatus handles PUT /trips/:id/status
func (c *TripController) UpdateTripStatus(ctx *gin.Context) {
	tripID := ctx.Param("id")

	var request dto.UpdateTripStatusRequest
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

	result, err := c.service.UpdateTripStatus(tripID, request)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Trip status updated successfully", result)
}
