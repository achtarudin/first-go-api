package http

import (
	"cutbray/first_api/domain/courier/entity"
	"cutbray/first_api/domain/courier/handler/request"
	"cutbray/first_api/domain/courier/usecase"
	"net/http"

	"cutbray/first_api/pkg/customerror"
	"cutbray/first_api/pkg/response"
	"cutbray/first_api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type courierHandler struct {
	router      *gin.RouterGroup
	middlewares []gin.HandlerFunc
	validator   *utils.Validator
	usecase     usecase.CourierUsecase
}

func NewCourierHandler(router *gin.RouterGroup, middlewares []gin.HandlerFunc, validator *utils.Validator, usecase usecase.CourierUsecase) *courierHandler {
	return &courierHandler{
		router:      router,
		middlewares: middlewares,
		validator:   validator,
		usecase:     usecase,
	}
}

func (h *courierHandler) RegisterRoute() {

	h.router.POST("/couriers/login", h.Login)
	h.router.POST("/couriers/register", h.Register)
	h.router.GET("/couriers/get-all", h.GetAllCouriers)
	h.router.GET("/couriers/get-by-long-lat", h.GetCourierByLongLat)
	h.router.GET("/couriers/find-nearest", h.FindNearestCourier)

	routeGroupMiddleware := h.router.Group("", h.middlewares...)
	routeGroupMiddleware.PUT("/couriers/update", h.Update)
	routeGroupMiddleware.DELETE("/couriers/delete", h.Delete)

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

	courier := json.ToCourierLogin()
	courierResult, err := h.usecase.Login(c.Request.Context(), courier.Email, courier.Password, utils.VerifyPassword)

	// If error occurs during usecase execution, return error response
	if err != nil {
		c.Error(err)
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

	courier := json.ToCourierRegister()
	createdCourier, err := h.usecase.Register(c.Request.Context(), &courier, utils.HashPassword)

	// If error occurs during usecase execution, return error response
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Register success",
		Data:    createdCourier,
	})
}

// GetAllCouriers godoc
//
//	@Summary	Get all couriers
//	@Tags		Couriers
//	@Accept		json
//	@Param		name		query	string	false	"search by name"
//	@Param		email		query	string	false	"search by email"
//	@Param		longitude	query	string	false	"search by longitude"				default(106.8260)
//	@Param		latitude	query	string	false	"search by latitude"				default(-6.1790)
//	@Param		radius		query	string	false	"radius in meter"					default(100)
//	@Param		per_page	query	int		false	"per page"							default(10)
//	@Param		page		query	int		false	"page"								default(1)
//	@Param		sort_by		query	string	false	"sort by (id, distance_in_meters)"	default(id)
//	@Param		order_by	query	string	false	"order by (ASC , DESC)"				default(ASC)
//	@Produce	json
//	@Success	200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success	201
//	@Failure	400
//	@Failure	404
//	@Failure	422
//	@Failure	500
//	@Router		/api/couriers/get-all [get]
func (h *courierHandler) GetAllCouriers(c *gin.Context) {

	var query request.GetAllCourierRequest
	err := c.ShouldBindQuery(&query)
	if err != nil {
		customErr := customerror.New(customerror.CodeInvalidInput, "Invalid input",
			nil, err,
		)
		c.Error(customErr)
		return
	}

	// Validate struct using validator
	errorMessage, isValid := h.validator.ValidateStruct(query)
	if isValid == false {
		customErr := customerror.NewFieldToAny(customerror.CodeValidationFailed, "Validation failed",
			errorMessage, nil,
		)
		c.Error(customErr)
		return
	}

	result, err := h.usecase.GetAllCouriers(c.Request.Context(), query.ToEntity())

	// If error occurs during usecase execution, return error response
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Get all couriers success",
		Data:    result,
	})
}

// GetCourierByLongLat godoc
//
//	@Summary	Get all couriers
//	@Tags		Couriers
//	@Accept		json
//	@Param		longitude	query	string	true	"search by longitude"				default(106.8260)
//	@Param		latitude	query	string	true	"search by latitude"				default(-6.1790)
//	@Param		per_page	query	int		false	"per page"							default(10)
//	@Param		page		query	int		false	"page"								default(1)
//	@Param		sort_by		query	string	false	"sort by (id, distance_in_meters)"	default(id)
//	@Param		order_by	query	string	false	"order by (ASC , DESC)"				default(ASC)
//	@Produce	json
//	@Success	200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success	201
//	@Failure	400
//	@Failure	404
//	@Failure	422
//	@Failure	500
//	@Router		/api/couriers/get-by-long-lat [get]
func (h *courierHandler) GetCourierByLongLat(c *gin.Context) {

	var query request.GetCourierByLongLatRequest
	err := c.ShouldBindQuery(&query)
	if err != nil {
		customErr := customerror.New(customerror.CodeInvalidInput, "Invalid input",
			nil, err,
		)
		c.Error(customErr)
		return
	}

	// Validate struct using validator
	errorMessage, isValid := h.validator.ValidateStruct(query)
	if isValid == false {
		customErr := customerror.NewFieldToAny(customerror.CodeValidationFailed, "Validation failed",
			errorMessage, nil,
		)
		c.Error(customErr)
		return
	}

	result, err := h.usecase.GetCourierByLongLat(c.Request.Context(), query.ToEntity())

	// If error occurs during usecase execution, return error response
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Get all couriers success",
		Data:    result,
	})
}

// FindNearestCourier godoc
//
//	@Summary	Get all couriers
//	@Tags		Couriers
//	@Accept		json
//	@Param		longitude	query	string	true	"search by longitude"	default(106.8260)
//	@Param		latitude	query	string	true	"search by latitude"	default(-6.1790)
//	@Param		radius		query	string	true	"radius in meter"		default(100)
//	@Produce	json
//	@Success	200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success	201
//	@Failure	400
//	@Failure	404
//	@Failure	422
//	@Failure	500
//	@Router		/api/couriers/find-nearest [get]
func (h *courierHandler) FindNearestCourier(c *gin.Context) {
	var query request.GetNearestRequest
	err := c.ShouldBindQuery(&query)

	if err != nil {
		customErr := customerror.New(customerror.CodeInvalidInput, "Invalid input",
			nil, err,
		)
		c.Error(customErr)
		return
	}

	// Validate struct using validator
	errorMessage, isValid := h.validator.ValidateStruct(query)
	if isValid == false {
		customErr := customerror.NewFieldToAny(customerror.CodeValidationFailed, "Validation failed",
			errorMessage, nil,
		)
		c.Error(customErr)
		return
	}

	result, err := h.usecase.GetNearestCourier(c.Request.Context(), query.ToEntity())

	// If error occurs during usecase execution, return error response
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Get all couriers success",
		Data:    map[string]*[]entity.Courier{"data": result},
	})
}

// Update godoc
//
//	@Security	ApiKeyAuth
//	@Summary	Update a courier
//	@Tags		Couriers
//	@Accept		json
//	@Param		payload	body	request.UpdateRequest	true	"json type"
//	@Produce	json
//	@Success	200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success	201
//	@Failure	400
//	@Failure	404
//	@Failure	422
//	@Failure	500
//	@Router		/api/couriers/update [put]
func (h *courierHandler) Update(c *gin.Context) {

	var json request.UpdateRequest
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

	userId, _ := c.Get("user_id")
	result, err := h.usecase.UpdateCourier(c.Request.Context(), int(userId.(float64)), json.ToEntity(), utils.HashPassword)

	// If error occurs during usecase execution, return error response
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Update courier success",
		Data:    result,
	})
}

// Delete godoc
//
//	@Security	ApiKeyAuth
//	@Summary	Delete a courier
//	@Tags		Couriers
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success	201
//	@Failure	400
//	@Failure	404
//	@Failure	422
//	@Failure	500
//	@Router		/api/couriers/delete [delete]
func (h *courierHandler) Delete(c *gin.Context) {
	userId, _ := c.Get("user_id")

	result, err := h.usecase.DeletedCourier(c.Request.Context(), int(userId.(float64)))

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Delete courier success",
		Data:    result,
	})
}
