package main

import (
	handler "cutbray/first_api/handler/http"
	"cutbray/first_api/infra"
	"cutbray/first_api/repository"
	"cutbray/first_api/utils"

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
	db, err := infra.NewDatabaseEnv()
	if err != nil {
		panic(err)
	}

	// Initialize validator
	validate := utils.NewValidator()

	// Initialize repositories
	repos := repository.NewRepositories(db.DB)
	_ = repos

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
	handler.NewAuthHandler(api, validate)
	handler.NewSwaggerHandler(server, "First Go API", "This is a sample server for the First Go API.")
	handler.NewHelloHandler(server)

	server.Run(":8080")
}
