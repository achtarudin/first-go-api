package main

import (
	"cutbray/first_api/infra"
	"cutbray/first_api/utils"

	swagger "cutbray/first_api/domain/docs/handler/http"
	hello "cutbray/first_api/domain/hello/handler/http"

	auth "cutbray/first_api/domain/auth/handler/http"

	"github.com/fatih/color"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {

	// Read environment variables
	viper.AutomaticEnv()

	// Initialize database
	_, err := infra.NewDatabaseEnv()
	if err != nil {
		panic(err)
	}

	// Initialize validator
	validate := utils.NewValidator()

	// Initialize Gin router
	if viper.GetBool("DEBUG") == false {
		color.Green("Service RUN on production mode")
		gin.SetMode(gin.ReleaseMode)
	}
	server := gin.Default()
	server.Use(cors.Default())
	server.SetTrustedProxies([]string{"127.0.0.1"})

	// API Router
	api := server.Group("/api")

	// Initialize handlers
	hello.NewHelloHandler(server)
	swagger.NewSwaggerHandler(server, "First GO API", "Documentation for First GO API")
	auth.NewAuthHandler(api, validate)

	server.Run(":8080")
}
