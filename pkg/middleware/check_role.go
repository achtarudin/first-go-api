package middleware

import (
	"cutbray/first_api/infra"
	"cutbray/first_api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CheckRoleRepository interface {
	IsCourier() gin.HandlerFunc
	IsMerchant() gin.HandlerFunc
}
type checkRoleRepository struct {
	db *infra.Database
}

func NewCheckRoleRepository(db *infra.Database) CheckRoleRepository {
	return &checkRoleRepository{
		db: db,
	}
}

func (r *checkRoleRepository) IsCourier() gin.HandlerFunc {
	return func(c *gin.Context) {
		if false {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.BindErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized",
			})
			return
		}
		c.Next()
	}
}

func (r *checkRoleRepository) IsMerchant() gin.HandlerFunc {
	return func(c *gin.Context) {
		if false {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.BindErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized",
			})
			return
		}
		c.Next()
	}
}
