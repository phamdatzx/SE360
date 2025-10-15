package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cấu trúc chung cho tất cả response
type ApiResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Hàm trả response thành công
func SuccessResponse(ctx *gin.Context, status int, message string, data interface{}) {
	ctx.JSON(http.StatusOK, ApiResponse{
		Status:  status,
		Message: message,
		Data:    data,
	})
}
