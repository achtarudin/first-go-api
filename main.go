package main

import (
	"cutbray/first_api/handler/http"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	http.NewHelloHandler(server)

	server.Run(":8080")
}
