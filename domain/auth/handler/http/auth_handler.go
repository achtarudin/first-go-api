package http

import (
	"cutbray/first_api/domain/auth/handler/middleware"
	"cutbray/first_api/domain/auth/handler/request"
	"cutbray/first_api/domain/auth/usecase"
	"cutbray/first_api/utils"
	"cutbray/first_api/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	router    *gin.RouterGroup
	usecase   usecase.AuthUsecase
	validator *utils.Validator
}

func NewAuthHandler(router *gin.RouterGroup, usecase usecase.AuthUsecase, validator *utils.Validator) {

	handler := &authHandler{
		router:    router,
		usecase:   usecase,
		validator: validator,
	}

	routeMiddleware := router.Use(middleware.JWTAuth())
	routeMiddleware.POST("/auth/login", handler.Login)
	routeMiddleware.POST("/auth/register", handler.Register)

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

	if errorMessage, isValid := h.validator.ValidateStruct(json); isValid == false {
		c.JSON(http.StatusUnprocessableEntity, response.BindErrorResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: "Validation failed",
			Errors:  errorMessage,
		})
		return
	}

	// Convert to entity and call usecase
	user := json.ToUserLogin()
	err := h.usecase.Login(c, &user, utils.VerifyPassword)

	// If error occurs during usecase execution, return error response
	if err != nil {
		c.JSON(http.StatusNotFound, response.BindErrorResponse{
			Status:  http.StatusNotFound,
			Message: "Not found",
			Errors: map[string]string{
				"error": err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Login success",
		Data:    user,
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

	// Validate input
	var json request.RegisterRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, response.BindErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Bad Request",
		})
		return
	}

	// Validate struct using validator
	if errorMessage, isValid := h.validator.ValidateStruct(json); isValid == false {
		c.JSON(http.StatusUnprocessableEntity, response.BindErrorResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: "Validation failed",
			Errors:  errorMessage,
		})
		return
	}

	// Convert to entity and call usecase
	user := json.ToUserRegister()
	err := h.usecase.Register(c, &user, utils.HashPassword)

	// If error occurs during usecase execution, return error response
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.BindErrorResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: "Validation failed",
			Errors: map[string]string{
				"error": err.Error(),
			},
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Register success",
		Data:    user,
	})
}
