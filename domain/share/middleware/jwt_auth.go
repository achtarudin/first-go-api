package middleware

import (
	"cutbray/first_api/utils"
	"cutbray/first_api/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.BindErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized",
				Errors: map[string]string{
					"error": "Unauthorized token not provided",
				},
			})
			return
		}

		// Verify token
		token, err := utils.VerifyToken(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.BindErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized",
				Errors: map[string]string{
					"error": err.Error(),
				},
			})
		}

		// Check token validity and extract user data
		valid, userData, err := utils.CheckTokenValid(token)
		if err != nil || valid == false {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.BindErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized",
				Errors: map[string]string{
					"error": err.Error(),
				},
			})
			return
		}

		// Store user data in context
		c.Set("user_id", userData["user_id"])
		c.Set("email", userData["email"])

		// Proceed to the next handler
		c.Next()
	}
}
