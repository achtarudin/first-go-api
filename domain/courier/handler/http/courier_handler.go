package http

import (
	"cutbray/first_api/domain/courier/handler/request"
	"cutbray/first_api/domain/courier/usecase"
	"net/http"

	"cutbray/first_api/pkg/response"
	"cutbray/first_api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type courierHandler struct {
	router    *gin.RouterGroup
	usecase   usecase.CourierUsecase
	validator *utils.Validator
}

func NewCourierHandler(router *gin.RouterGroup, usecase usecase.CourierUsecase, validator *utils.Validator) *courierHandler {
	return &courierHandler{
		router:    router,
		usecase:   usecase,
		validator: validator,
	}
}

func (h *courierHandler) RegisterRoute() {
	h.router.POST("/couriers/login", h.Login)
	h.router.POST("/couriers/register", h.Register)
}

// Login godoc
//
//	@Summary	Authenticate couriers with email and password
//	@Tags		Couriers
//	@Accept		json
//	@Param		payload	body	request.LoginRequest	true	"json type"
//	@Produce	json
//	@Success	200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success	201
//	@Failure	400
//	@Failure	404
//	@Failure	422
//	@Failure	500
//	@Router		/api/couriers/login [post]
func (h *courierHandler) Login(c *gin.Context) {
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

	courier := json.ToCourierLogin()
	courierResult, err := h.usecase.Login(c, &courier, utils.VerifyPassword)

	// If error occurs during usecase execution, return error response
	if err != nil {
		c.JSON(http.StatusNotFound, response.BindErrorResponse{
			Status:  http.StatusNotFound,
			Message: "Not Found",
			Errors: map[string]string{
				"error": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Login success",
		Data:    courierResult,
	})
}

// Register godoc
//
//	@Summary	Register a new courier
//	@Tags		Couriers
//	@Accept		json
//	@Param		payload	body	request.RegisterRequest	true	"json type"
//	@Produce	json
//	@Success	200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success	201
//	@Failure	400
//	@Failure	404
//	@Failure	422
//	@Failure	500
//	@Router		/api/couriers/register [post]
func (h *courierHandler) Register(c *gin.Context) {
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

	courier := json.ToCourierRegister()
	createdCourier, err := h.usecase.Register(c, &courier, utils.HashPassword)

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
	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Register success",
		Data:    createdCourier,
	})
}
