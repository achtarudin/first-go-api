package http

import (
	"cutbray/first_api/domain/merchant/handler/request"
	"cutbray/first_api/domain/merchant/usecase"
	"net/http"

	"cutbray/first_api/pkg/customerror"
	"cutbray/first_api/pkg/response"
	"cutbray/first_api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type merchantHandler struct {
	router      *gin.RouterGroup
	middlewares []gin.HandlerFunc
	validator   *utils.Validator
	usecase     usecase.MerchantUsecase
}

func NewMerchantHandler(
	router *gin.RouterGroup, middlewares []gin.HandlerFunc,
	validator *utils.Validator, usecase usecase.MerchantUsecase) *merchantHandler {
	return &merchantHandler{
		router:      router,
		middlewares: middlewares,
		validator:   validator,
		usecase:     usecase,
	}
}

func (h *merchantHandler) RegisterRoute() {
	h.router.POST("/merchants/login", h.Login)
	h.router.POST("/merchants/register", h.Register)
	h.router.GET("/merchants/get-all", h.GetAll)
	h.router.GET("/merchants/get-by-id", h.GetById)
	h.router.GET("/merchants/get-by-user-id", h.GetByUserId)

	routeGroupMiddleware := h.router.Group("", h.middlewares...)
	routeGroupMiddleware.PUT("/merchants/update", h.Update)
	routeGroupMiddleware.DELETE("/merchants/delete", h.Delete)
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
	err := c.ShouldBindJSON(&json)
	if err != nil {
		customErr := customerror.New(customerror.CodeInvalidInput, "Invalid input",
			nil, err,
		)
		c.Error(customErr)
		return
	}

	errorMessage, isValid := h.validator.ValidateStruct(json)
	if isValid == false {
		customErr := customerror.NewFieldToAny(customerror.CodeValidationFailed, "Validation failed",
			errorMessage, nil,
		)
		c.Error(customErr)
		return
	}

	result, err := h.usecase.Login(c.Request.Context(), json.ToMerchantLogin(), utils.VerifyPassword, utils.GenerateTokenFromIdAndEmail)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Login success",
		Data:    result,
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
	err := c.ShouldBindJSON(&json)
	if err != nil {
		customErr := customerror.New(customerror.CodeInvalidInput, "Invalid input",
			nil, err,
		)
		c.Error(customErr)
		return
	}

	errorMessage, isValid := h.validator.ValidateStruct(json)
	if isValid == false {
		customErr := customerror.NewFieldToAny(customerror.CodeValidationFailed, "Validation failed",
			errorMessage, nil,
		)
		c.Error(customErr)
		return
	}

	merchant, err := h.usecase.Register(c.Request.Context(), json.ToMerchantRegister(), utils.HashPassword)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Register success",
		Data:    merchant,
	})
}

// GetAll godoc
//
//	@Summary	Get all merchants
//	@Tags		Merchants
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success	201
//	@Failure	400
//	@Failure	404
//	@Failure	422
//	@Failure	500
//	@Router		/api/merchants/get-all [get]
func (h *merchantHandler) GetAll(c *gin.Context) {

	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Get all merchants success",
	})
}

// GetById godoc
//
//	@Summary	Get a merchant by ID
//	@Tags		Merchants
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success	201
//	@Failure	400
//	@Failure	404
//	@Failure	422
//	@Failure	500
//	@Router		/api/merchants/get-by-id [get]
func (h *merchantHandler) GetById(c *gin.Context) {

	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Get merchant by ID success",
	})
}

// GetByUserId godoc
//
//	@Summary	Get merchants by User ID
//	@Tags		Merchants
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success	201
//	@Failure	400
//	@Failure	404
//	@Failure	422
//	@Failure	500
//	@Router		/api/merchants/get-by-user-id [get]
func (h *merchantHandler) GetByUserId(c *gin.Context) {
	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Get merchants by User ID success",
	})
}

// Update godoc
//
//	@Security	ApiKeyAuth
//	@Summary	Update a merchant
//	@Tags		Merchants
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success	201
//	@Failure	400
//	@Failure	404
//	@Failure	422
//	@Failure	500
//	@Router		/api/merchants/update [put]
func (h *merchantHandler) Update(c *gin.Context) {
	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Update merchants by User ID success",
	})
}

// Delete godoc
//
//	@Security	ApiKeyAuth
//	@Summary	Delete a merchant
//	@Tags		Merchants
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success	201
//	@Failure	400
//	@Failure	404
//	@Failure	422
//	@Failure	500
//	@Router		/api/merchants/delete [delete]
func (h *merchantHandler) Delete(c *gin.Context) {
	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Update merchants by User ID success",
	})
}
