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

// Hello godoc
// @Summary      Show a hello message
// @Description  Hello Handler
// @Tags         Hello
// @Accept       json
// @Produce      json
// @Success      200
// @Success      201
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       / [get]
func (h *helloHandler) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, response.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Hello, Ngopi yuk!",
		Data:    []any{},
	})
}
