package main

import (
	"cutbray/first_api/infra"
	"cutbray/first_api/pkg/middleware"
	"cutbray/first_api/pkg/utils"
	"fmt"

	swagger "cutbray/first_api/domain/docs/handler/http"
	hello "cutbray/first_api/domain/hello/handler/http"

	auth "cutbray/first_api/domain/auth/handler/http"
	repoAuth "cutbray/first_api/domain/auth/repository"
	usecaseAuth "cutbray/first_api/domain/auth/usecase"

	"github.com/fatih/color"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {

	// Load configuration
	configViper, err := loadConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println(configViper.GetString("DB_HOST"))
	// Initialize database
	db, err := infra.NewDatabase(infra.DatabaseConfig{
		Host:     configViper.GetString("DB_HOST"),
		Port:     configViper.GetInt("DB_PORT"),
		User:     configViper.GetString("DB_USER"),
		Password: configViper.GetString("DB_PASSWORD"),
		DBName:   configViper.GetString("DB_DATABASE"),
	})
	if err != nil {
		panic(err)
	}

	// Initialize validator
	validate := utils.NewValidator()

	// Initialize Gin router
	if configViper.GetBool("DEBUG") == false {
		color.Green("Service RUN on production mode")
		gin.SetMode(gin.ReleaseMode)
	}

	server := gin.Default()
	server.Use(cors.Default())
	server.SetTrustedProxies([]string{"127.0.0.1"})

	// Initialize Repository
	repoAuth := repoAuth.NewAuthRepository(db.DB)

	// Initialize Usecase
	usecaseAuth := usecaseAuth.NewAuthUsecase(repoAuth)

	// Initialize Middleware
	authMiddleware := middleware.JWTAuth()

	// Initialize API group
	api := server.Group("/api")

	// Initialize hello
	{
		helloHandler := hello.NewHelloHandler(server, &authMiddleware)
		helloHandler.RegisterRoute()
	}

	// Initialize swagger
	{
		swaggerHadler := swagger.NewSwaggerHandler(server, "First GO API", "Documentation for First GO API")
		swaggerHadler.RegisterRoute()
	}

	// Initialize auth
	{
		authHandler := auth.NewAuthHandler(api, usecaseAuth, validate)
		authHandler.RegisterRoute()
	}

	server.Run(":8080")
}

func loadConfig() (*viper.Viper, error) {
	config := infra.NewAppConfig()

	err := config.LoadEnvConfig(nil)
	if err != nil {
		return nil, err
	}

	err = config.LoadTranslationConfig(nil)
	if err != nil {
		return nil, err
	}

	return config.GetViper(), nil
}
