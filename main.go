package main

import (
	handler "cutbray/first_api/handler/http"

	"github.com/fatih/color"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	// Read environment variables
	viper.AutomaticEnv()

	// Set gin mode
	if viper.GetBool("DEBUG") == false {
		color.Green("Service RUN on production mode")
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	server := gin.Default()
	server.Use(cors.Default())
	server.SetTrustedProxies([]string{"127.0.0.1"})

	// Initialize handlers
	handler.NewSwaggerHandler(server, "First Go API", "This is a sample server for the First Go API.")
	handler.NewHelloHandler(server)

	server.Run(":8080")
}
