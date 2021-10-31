package router

import (
	"github.com/RobotApocalypse/configuration"
	"github.com/RobotApocalypse/constants"
	"github.com/RobotApocalypse/controllers"
	"github.com/RobotApocalypse/docs"
	"github.com/RobotApocalypse/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(l *logrus.Entry, cfg configuration.Config) *gin.Engine  {

	router := gin.Default()


	newMiddleware := middleware.NewMiddleware(l)

	router.Use(newMiddleware.RecoverHandler())

	docs.SwaggerInfo.Title = cfg.GetString(constants.ServiceName)
	docs.SwaggerInfo.Description = "Lists down the API to tackle ROBOT Apocalypse"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	healthCheckCtl := controllers.NewHealthCheckCtrl(l)
	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/api/rob/v1/healthz", healthCheckCtl.HealthCheck)

	return router
}

