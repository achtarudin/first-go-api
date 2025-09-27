package middleware

import (
	"cutbray/first_api/pkg/customerror"
	"cutbray/first_api/pkg/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Step1: Process the request first.

		// Step2: Check if any errors were added to the context
		if len(c.Errors) > 0 {
			// Step3: Use the last error
			err := c.Errors.Last().Err

			// Step4: Respond with a generic error message
			var customError *customerror.CustomError
			if ok := errors.As(err, &customError); ok {
				c.JSON(customError.StatusCode(), response.BindErrorResponse{
					Status:  customError.StatusCode(),
					Message: customError.Message(),
					Errors:  customError.Fields(),
				})

				return
			}

			c.JSON(http.StatusInternalServerError, response.BindErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
		}

		// Any other steps if no errors are found
	}
}
