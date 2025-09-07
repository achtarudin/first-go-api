package main

import (
	handler "cutbray/first_api/handler/http"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		c := color.New(color.BgRed, color.FgWhite).Add(color.Bold)
		panic(c.Sprintf("Fatal error config file: %v", err))
	}

	// Read environment variables
	viper.AutomaticEnv()

	// Set gin mode
	if viper.GetBool("DEBUG") == false {
		color.Green("Service RUN on production mode")
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	server := gin.Default()
	server.SetTrustedProxies([]string{"127.0.0.1"})

	// Initialize handlers
	handler.NewHelloHandler(server)

	server.Run(":8080")
}
