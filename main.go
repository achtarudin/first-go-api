package main

import (
	"cutbray/first_api/infra"
	"cutbray/first_api/pkg/middleware"
	"cutbray/first_api/pkg/utils"

	swagger "cutbray/first_api/domain/docs/handler/http"
	hello "cutbray/first_api/domain/hello/handler/http"

	auth "cutbray/first_api/domain/auth/handler/http"
	authRepo "cutbray/first_api/domain/auth/repository"
	authUsecase "cutbray/first_api/domain/auth/usecase"

	courier "cutbray/first_api/domain/courier/handler/http"
	courierRepo "cutbray/first_api/domain/courier/repository"
	courierUsecase "cutbray/first_api/domain/courier/usecase"

	merchant "cutbray/first_api/domain/merchant/handler/http"
	merchantRepo "cutbray/first_api/domain/merchant/repository"
	merchantUsecase "cutbray/first_api/domain/merchant/usecase"

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

	// Initialize swagger
	{
		swaggerHadler := swagger.NewSwaggerHandler(server, "First GO API", "Documentation for First GO API")
		swaggerHadler.RegisterRoute()
	}

	{
		// Initialize Middleware
		authMiddleware := middleware.JWTAuth()

		// Initialize hello
		{
			helloHandler := hello.NewHelloHandler(server, &authMiddleware)
			helloHandler.RegisterRoute()
		}
	}

	{
		// Initialize API group
		api := server.Group("/api")
		checkRoleMiddleware := middleware.NewCheckRoleRepository(db)
		api.Use(checkRoleMiddleware.IsMerchant())
		api.Use(checkRoleMiddleware.IsCourier())
		// Initialize auth
		{
			repoAuth := authRepo.NewAuthRepository(db)
			usecaseAuth := authUsecase.NewAuthUsecase(repoAuth)
			authHandler := auth.NewAuthHandler(api, usecaseAuth, validate)
			authHandler.RegisterRoute()
		}

		// Initialize courier
		{
			courierRepo := courierRepo.NewCourierRepository(db)
			usecaseCourier := courierUsecase.NewCourierUsecase(courierRepo)
			courierHandler := courier.NewCourierHandler(api, usecaseCourier, validate)
			courierHandler.RegisterRoute()
		}

		// Initialize merchant
		{
			merchantRepo := merchantRepo.NewMerchantRepository(db)
			usecaseMerchant := merchantUsecase.NewMerchantUsecase(merchantRepo)
			merchantHandler := merchant.NewMerchantHandler(api, usecaseMerchant, validate)
			merchantHandler.RegisterRoute()
		}
	}

	server.Run(":8080")
}

func loadConfig() (*viper.Viper, error) {
	config := infra.NewAppConfig()

	err := config.LoadEnvConfig(nil)
	if err != nil {
		return nil, err
	}

	// err = config.LoadTranslationConfig(nil)
	// if err != nil {
	// 	return nil, err
	// }

	return config.GetViper(), nil
}
