package http

import (
	"cutbray/first_api/domain/auth/handler/request"
	"cutbray/first_api/utils"
	"cutbray/first_api/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	router    *gin.RouterGroup
	validator *utils.Validator
}

func NewAuthHandler(router *gin.RouterGroup, validator *utils.Validator) {

	handler := &authHandler{
		router:    router,
		validator: validator,
	}

	router.POST("/auth/login", handler.Login)
	router.POST("/auth/register", handler.Register)

}

// Login godoc
//
//	@Description	Authenticate user with email and password
//	@Tags			Auth
//	@Accept			json
//	@Param			payload	body	request.LoginRequest	true	"json type"
//	@Produce		json
//	@Success		200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success		201
//	@Failure		400
//	@Failure		404
//	@Failure		422
//	@Failure		500
//	@Router			/api/auth/login [post]
func (h *authHandler) Login(c *gin.Context) {

	var json request.LoginRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, response.BindErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Bad Request",
		})
		return
	}

	if err := h.validator.ValidateStruct(json); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.BindErrorResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: "Validation failed",
			Errors:  err,
		})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Hello, Ngopi yuk!",
		Data:    []any{},
	})
}

// Register godoc
//
//	@Description	Register a new user
//	@Tags			Auth
//	@Accept			json
//	@Param			payload	body	request.RegisterRequest	true	"json type"
//	@Produce		json
//	@Success		200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success		201
//	@Failure		400
//	@Failure		404
//	@Failure		422
//	@Failure		500
//	@Router			/api/auth/register [post]
func (h *authHandler) Register(c *gin.Context) {

	var json request.RegisterRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, response.BindErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Bad Request",
		})
		return
	}

	if err := h.validator.ValidateStruct(json); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.BindErrorResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: "Validation failed",
			Errors:  err,
		})
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Hello, Ngopi yuk!",
		Data:    []any{},
	})
}
