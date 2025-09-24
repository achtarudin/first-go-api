package middleware

import (
	"cutbray/first_api/infra"
	"cutbray/first_api/pkg/model"
	"cutbray/first_api/pkg/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

		userId, existsId := c.Get("user_id")
		email, existsEmail := c.Get("email")

		_ = userId
		if !existsId || !existsEmail {
			c.AbortWithStatusJSON(http.StatusForbidden, response.BindErrorResponse{
				Status:  http.StatusForbidden,
				Message: "Forbidden",
			})
			return
		}

		user := model.User{}

		result := r.db.DB.
			Where("email = ?", email).
			Where("id = ?", userId).
			Where("id IN (?)", r.db.DB.Model(&model.Courier{}).Select("user_id")).
			First(&user)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.AbortWithStatusJSON(http.StatusForbidden, response.BindErrorResponse{
					Status:  http.StatusForbidden,
					Message: "Forbidden",
					Errors:  result.Error.Error(),
				})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.BindErrorResponse{
					Status:  http.StatusInternalServerError,
					Message: "Internal Server Error",
				})
			}

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
