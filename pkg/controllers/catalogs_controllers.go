package controllers

import (
	"github.com/afiskon/promtail-client/promtail"
	"github.com/gin-gonic/gin"
	"github.com/rusrafkasimov/catalogs/pkg/usecases"
)

type CatalogsController struct {
	logger     promtail.Client
	catalogsUC usecases.CatalogsUseCase
}

func NewCatalogsController(uc usecases.CatalogsUseCase, logger promtail.Client) *CatalogsController {
	return &CatalogsController{
		logger:     logger,
		catalogsUC: uc,
	}
}

func (cc *CatalogsController) CreateCatalog(c *gin.Context) {

}

func (cc *CatalogsController) GetCatalogs(c *gin.Context) {

}

func (cc *CatalogsController) GetCatalogCategories(c *gin.Context) {

}

func (cc *CatalogsController) GetCatalogByID(c *gin.Context) {

}

func (cc *CatalogsController) UpdateCatalogByID(c *gin.Context) {

}

func (cc *CatalogsController) DeleteCatalogByID(c *gin.Context) {

}
