package http

import (
	"cutbray/first_api/domain/auth/entity"
	"cutbray/first_api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type helloHandler struct {
	server     *gin.Engine
	middleware *gin.HandlerFunc
}

func NewHelloHandler(server *gin.Engine, middleware *gin.HandlerFunc) *helloHandler {

	return &helloHandler{
		server:     server,
		middleware: middleware,
	}
}

func (h *helloHandler) RegisterRoute() {
	h.server.GET("/", *h.middleware, h.Hello)
	h.server.GET("/health", h.Health)

	h.server.POST("/post-hello", h.PostBodyHello)
	h.server.POST("/post-hello-form", h.PostFormDataHello)
}

// Hello godoc
//
//	@Description	Show a hello message
//	@Security		ApiKeyAuth
//	@Tags			Hello
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string									false	"search by name"
//	@Success		200		{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success		201
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/ [get]
func (h *helloHandler) Hello(c *gin.Context) {

	userId, _ := c.Get("user_id")
	email, _ := c.Get("email")

	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Hello, Ngopi yuk!",
		Data:    entity.User{ID: int(userId.(float64)), Email: email.(string)},
	})
}

// Health godoc
//
//	@Description	Check the health of the application
//	@Tags			Hello
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success		201
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/health [get]
func (h *helloHandler) Health(c *gin.Context) {

	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "The application is healthy",
	})
}

// PostBodyHello godoc
//
//	@Description	Post a hello message using json body
//	@Tags			Hello
//	@Accept			json
//	@Param			payload	body	request.HelloRequest	true	"json type"
//	@Produce		json
//	@Success		200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success		201
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/post-hello [post]
func (h *helloHandler) PostBodyHello(c *gin.Context) {
	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Hello, Ngopi yuk!",
		Data:    []any{},
	})
}

// PostFormDataHello godoc
//
//	@Description	Post a hello message using form data
//	@Tags			Hello
//	@Accept			multipart/form-data
//	@Param			name	formData	string	true	"Name"
//	@Param			age		formData	int		true	"Age"
//	@Param			email	formData	string	true	"Email"
//	@Param			files	formData	file	true	"File to upload"
//	@Produce		json
//	@Success		200	{object}	response.SuccessResponse{data=[]any}	"success response so the data field is array of any type"
//	@Success		201
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/post-hello-form [post]
func (h *helloHandler) PostFormDataHello(c *gin.Context) {
	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Hello, Ngopi yuk!",
		Data:    []any{},
	})
}
