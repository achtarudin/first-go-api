package http

import (
	"cutbray/first_api/handler/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type helloHandler struct {
	server *gin.Engine
}

func NewHelloHandler(server *gin.Engine) {

	handler := &helloHandler{
		server: server,
	}

	server.GET("/", handler.Hello)

}

func (h *helloHandler) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Hello, asas",
		Data:    []any{},
	})
}
