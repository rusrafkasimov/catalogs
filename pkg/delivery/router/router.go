package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MapUrl(router *gin.Engine, appCtx *ApplicationContext) {
	router.Use(SetMiddleware())

	// Set Middleware
	authorized := router.Group("/")
	authorized.Use(SetMiddlewareAuthentication())

	// Base routes
	authorized.OPTIONS("/catalog", appCtx.CatalogsController.CreateCatalog)
	authorized.POST("/catalog", appCtx.CatalogsController.CreateCatalog)
	authorized.GET("/catalog", appCtx.CatalogsController.GetCatalogs)
	authorized.PUT("/catalog", appCtx.CatalogsController.UpdateCatalog)
	authorized.OPTIONS("/catalog/:id", appCtx.CatalogsController.GetCatalogByID)
	authorized.GET("/catalog/:id", appCtx.CatalogsController.GetCatalogByID)
	authorized.DELETE("/catalog/:id", appCtx.CatalogsController.DeleteCatalog)
	authorized.OPTIONS("/categories", appCtx.CatalogsController.GetCatalogCategories)
	authorized.GET("/categories", appCtx.CatalogsController.GetCatalogCategories)


	// System Routes
	router.GET("/ping", PingHandler)
	router.GET("/health", HealthHandler)

	// Swagger Route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}