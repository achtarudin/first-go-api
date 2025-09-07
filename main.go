package main

import (
	"cutbray/first_api/docs"
	handler "cutbray/first_api/handler/http"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	// Read environment variables
	viper.AutomaticEnv()

	// Set gin mode
	if viper.GetBool("DEBUG") == false {
		color.Green("Service RUN on production mode")
		gin.SetMode(gin.ReleaseMode)
	}

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "First Go API"
	docs.SwaggerInfo.Description = "This is a sample server for the First Go API."
	docs.SwaggerInfo.Version = "1.0"
	// docs.SwaggerInfo.Host = "localhost:8080"
	// docs.SwaggerInfo.BasePath = "/"
	// docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Initialize Gin router
	server := gin.Default()
	server.SetTrustedProxies([]string{"127.0.0.1"})

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize handlers
	handler.NewHelloHandler(server)

	server.Run(":8080")
}
