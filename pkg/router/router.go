package router

import (
	"net/http"
	_ "qraven/cmd/api/docs"
	"qraven/internal/config"
	"qraven/pkg/middleware"
	"qraven/pkg/repository/storage"
	"qraven/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(logger *utils.Logger, validator *validator.Validate, db *storage.Database, appConfiguration *config.App) *gin.Engine {
	if appConfiguration.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	cspMiddleware := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline';")
			c.Next()
		}
	}

	// docs.SwaggerInfo.Title = "Swagger Example API"
	// docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	// docs.SwaggerInfo.Version = "1.0"
	// docs.SwaggerInfo.Host = "petstore.swagger.io"
	// docs.SwaggerInfo.BasePath = "/v2"
	// docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r := gin.New()

	// Middlewares
	// r.Use(gin.Logger())
	r.ForwardedByClientIP = true
	r.SetTrustedProxies(config.GetConfig().Server.TrustedProxies)
	r.Use(middleware.Security())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.Metrics(config.GetConfig()))
	// r.Use(middleware.GzipWithExclusion("/metrics"))
	r.MaxMultipartMemory = 3 << 20
	r.Use(cspMiddleware())

	// routers
	ApiVersion := "api/v1"
	api := r.Group(ApiVersion)
	Auth(r, ApiVersion, validator, db, logger)
	Event(r, ApiVersion, validator, db, logger)
	Ticket(r, ApiVersion, validator, db, logger)

	r.StaticFile("/swagger.yaml", "static/swagger.yaml")
	url := ginSwagger.URL("/swagger.yaml")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 200,
			"message":     "Welcome to QRaven API",
			"status":      http.StatusOK,
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"name":        "Not Found",
			"message":     "Page not found.",
			"status_code": 404,
			"status":      http.StatusNotFound,
		})
	})

	// Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}
