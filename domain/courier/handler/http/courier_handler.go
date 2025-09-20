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
	h.router.GET("/couriers/get-all", h.GetAllCouriers)
	h.router.GET("/couriers/get-by-long-lat", h.GetCourierByLongLat)
	h.router.GET("/couriers/find-nearest", h.FindNearestCourier)
	h.router.PUT("/couriers/update", h.Update)
	h.router.DELETE("/couriers/delete", h.Delete)

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

// GetAllCouriers godoc
//
//	@Summary	Get all couriers
//	@Tags		Couriers
//	@Accept		json
//	@Param		name		query	string	false	"search by name"
//	@Param		email		query	string	false	"search by email"
//	@Param		longitude	query	string	false	"search by longitude"				default(106.8260)
//	@Param		latitude	query	string	false	"search by latitude"				default(-6.1790)
//	@Param		radius		query	int		false	"radius in meter"					default(0)
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
	// Bind query parameters to struct
	err := c.ShouldBindQuery(&query)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.BindErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Bad Request",
		})

		return
	}

	// Validate struct using validator
	errorMessage, isValid := h.validator.ValidateStruct(query)
	if isValid == false {
		c.JSON(http.StatusBadRequest, response.BindErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Bad Request",
			Errors:  errorMessage,
		})
		return
	}

	result, err := h.usecase.GetAllCouriers(c, query.ToEntity())

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.BindErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal Server Error",
			Errors: map[string]string{
				"error": err.Error(),
			},
		})
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
//	@Param		longitude	query	string	false	"search by longitude"
//	@Param		latitude	query	string	false	"search by latitude"
//	@Produce	json
//	@Success	200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success	201
//	@Failure	400
//	@Failure	404
//	@Failure	422
//	@Failure	500
//	@Router		/api/couriers/get-by-long-lat [get]
func (h *courierHandler) GetCourierByLongLat(c *gin.Context) {
	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Get courier by longitude and latitude success",
	})
}

// FindNearestCourier godoc
//
//	@Summary	Get all couriers
//	@Tags		Couriers
//	@Accept		json
//	@Param		latitude	query	string	false	"search by latitude"
//	@Param		longitude	query	string	false	"search by longitude"
//	@Produce	json
//	@Success	200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success	201
//	@Failure	400
//	@Failure	404
//	@Failure	422
//	@Failure	500
//	@Router		/api/couriers/find-nearest [get]
func (h *courierHandler) FindNearestCourier(c *gin.Context) {
	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Get nearest courier success",
	})
}

// Update godoc
//
//	@Summary	Update a courier
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
//	@Router		/api/couriers/update [put]
func (h *courierHandler) Update(c *gin.Context) {
	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Update courier success",
	})
}

// Update godoc
//
//	@Summary	Update a courier
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
//	@Router		/api/couriers/delete [delete]
func (h *courierHandler) Delete(c *gin.Context) {
	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Delete courier success",
	})
}
