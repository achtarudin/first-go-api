package http

import (
	"cutbray/first_api/domain/merchant/handler/request"
	"cutbray/first_api/domain/merchant/usecase"
	"net/http"

	"cutbray/first_api/pkg/response"
	"cutbray/first_api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type merchantHandler struct {
	router    *gin.RouterGroup
	usecase   usecase.MerchantUsecase
	validator *utils.Validator
}

func NewMerchantHandler(router *gin.RouterGroup, usecase usecase.MerchantUsecase, validator *utils.Validator) *merchantHandler {
	return &merchantHandler{
		router:    router,
		usecase:   usecase,
		validator: validator,
	}
}

func (h *merchantHandler) RegisterRoute() {
	h.router.POST("/merchants/login", h.Login)
	h.router.POST("/merchants/register", h.Register)
}

// Login godoc
//
//	@Summary	Authenticate merchants with email and password
//	@Tags		Merchants
//	@Accept		json
//	@Param		payload	body	request.LoginRequest	true	"json type"
//	@Produce	json
//	@Success	200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success	201
//	@Failure	400
//	@Failure	404
//	@Failure	422
//	@Failure	500
//	@Router		/api/merchants/login [post]
func (h *merchantHandler) Login(c *gin.Context) {
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

	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Login success",
		Data:    []any{},
	})
}

// Register godoc
//
//	@Summary	Register a new merchant
//	@Tags		Merchants
//	@Accept		json
//	@Param		payload	body	request.RegisterRequest	true	"json type"
//	@Produce	json
//	@Success	200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success	201
//	@Failure	400
//	@Failure	404
//	@Failure	422
//	@Failure	500
//	@Router		/api/merchants/register [post]
func (h *merchantHandler) Register(c *gin.Context) {
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

	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Register success",
		Data:    []any{},
	})
}
