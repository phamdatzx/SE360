package middleware

import (
	"errors"
	"log"
	"net/http"

	"user-service/error"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Check if it's our custom AppError
			var appErr *error.AppError
			if errors.As(err, &appErr) {
				// Log internal errors
				if appErr.Err != nil {
					log.Printf("Error [%s]: %v", appErr.Code, appErr.Err)
				}

				c.JSON(appErr.Code, gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
				})
				return
			}

			// Handle unknown errors
			log.Printf("Unexpected error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			})
		}
	}
}
