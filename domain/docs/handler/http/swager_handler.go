package http

import (
	"cutbray/first_api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
)

type swaggerHandler struct {
	server *gin.Engine
}

func NewSwaggerHandler(server *gin.Engine, title string, description string) *swaggerHandler {
	docs.SwaggerInfo.Title = title
	docs.SwaggerInfo.Description = description
	docs.SwaggerInfo.Version = "1.0"
	// docs.SwaggerInfo.Host = "localhost:8080"
	// docs.SwaggerInfo.BasePath = "/"
	// docs.SwaggerInfo.Schemes = []string{"http", "https"}

	server.GET("/swagger/*any", ginSwagger.CustomWrapHandler(&ginSwagger.Config{
		URL:                      "doc.json",
		DocExpansion:             "list",
		InstanceName:             swag.Name,
		Title:                    title,
		DefaultModelsExpandDepth: -1,
		DeepLinking:              true,
		PersistAuthorization:     true,
		Oauth2DefaultClientID:    "",
		Oauth2UsePkce:            true,
	}, swaggerFiles.Handler))

	return &swaggerHandler{
		server: server,
	}
}
