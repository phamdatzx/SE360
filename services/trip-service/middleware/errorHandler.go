package middleware

import (
	apperror "trip-service/error"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Check if it's our custom AppError
			if appErr, ok := err.(*apperror.AppError); ok {
				c.JSON(appErr.StatusCode, gin.H{
					"error": appErr.Message,
				})
				return
			}

			// Default error response
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		}
	}
}
