package http

import (
	"cutbray/first_api/domain/auth/handler/request"
	"cutbray/first_api/domain/auth/usecase"
	"cutbray/first_api/pkg/customerror"
	"cutbray/first_api/pkg/response"
	"cutbray/first_api/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	router    *gin.RouterGroup
	usecase   usecase.AuthUsecase
	validator *utils.Validator
}

func NewAuthHandler(router *gin.RouterGroup, usecase usecase.AuthUsecase, validator *utils.Validator) *authHandler {
	return &authHandler{
		router:    router,
		usecase:   usecase,
		validator: validator,
	}
}

func (h *authHandler) RegisterRoute() {
	h.router.POST("/auth/login", h.Login)
	h.router.POST("/auth/register", h.Register)
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
	err := h.usecase.Login(c, &user, utils.VerifyPassword, utils.GenerateToken)

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
	// Return success response
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
	err := c.ShouldBindJSON(&json)
	if err != nil {
		customErr := customerror.New(customerror.CodeInvalidInput, "Invalid input",
			nil, err,
		)
		c.Error(customErr)
		return
	}

	// Validate struct using validator
	errorMessage, isValid := h.validator.ValidateStruct(json)
	if isValid == false {
		customErr := customerror.NewFieldToAny(customerror.CodeValidationFailed, "Validation failed",
			errorMessage, nil,
		)
		c.Error(customErr)
		return
	}

	// Convert to entity and call usecase
	user := json.ToUserRegister()
	err = h.usecase.Register(c, &user, utils.HashPassword)

	// If error occurs during usecase execution, return error response
	if err != nil {
		c.Error(err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Register success",
		Data:    user,
	})
}
