package middleware

import (
	"cutbray/first_api/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		isTokenValid := true

		if authHeader == "" && !isTokenValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.BindErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized",
			})
			return
		}

		c.Next()
	}
}

// package middleware

// import (
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt/v5"
// 	"cutbray/first_api/utils"
// )

// func JWTAuth() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Ambil token dari header Authorization: Bearer <token>
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
// 			return
// 		}

// 		// Split "Bearer <token>"
// 		parts := strings.Split(authHeader, " ")
// 		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
// 			return
// 		}
// 		tokenString := parts[1]

// 		// Verify token
// 		token, err := utils.VerifyToken(tokenString)
// 		if err != nil || !token.Valid {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
// 			return
// 		}

// 		// Ambil claims
// 		claims, ok := token.Claims.(jwt.MapClaims)
// 		if !ok {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
// 			return
// 		}

// 		// Check exp & iat
// 		now := time.Now().Unix()
// 		exp, okExp := claims["exp"].(float64)
// 		iat, okIat := claims["iat"].(float64)
// 		if !okExp || !okIat || int64(exp) < now || int64(iat) > now {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
// 			return
// 		}

// 		// Set user_id ke context
// 		c.Set("user_id", claims["user_id"])
// 		c.Next()
// 	}
// }
