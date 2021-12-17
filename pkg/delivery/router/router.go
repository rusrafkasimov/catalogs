package router

import (
	"github.com/gin-gonic/gin"
)

func MapUrl(router *gin.Engine, appCtx *ApplicationContext) {
	router.Use(SetMiddleware())

	// Set Middleware
	authorized := router.Group("/")
	authorized.Use(SetMiddlewareAuthentication())

	// Base routes

	// System Routes
	router.GET("/ping", PingHandler)
	router.GET("/health", HealthHandler)

	// Swagger Route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}