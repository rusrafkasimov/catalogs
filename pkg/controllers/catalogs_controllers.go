package controllers

import (
	"context"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/rusrafkasimov/catalogs/internal/errs"
	"github.com/rusrafkasimov/catalogs/internal/trace"
	"github.com/rusrafkasimov/catalogs/pkg/dto"
	"github.com/rusrafkasimov/catalogs/pkg/usecases"
	"net/http"
	"time"
)

const (
	errInvID      = "empty or invalid id parameter"
	errInvJSON    = "invalid json body"
	serverTimeout = 10 * time.Second
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

// CreateCatalog godoc
// @Summary Create catalog
// @Description Get JSON CatalogRequest, return JSON CreateCatalogResponse
// @Tags Catalog
// @Produce  json
// @Content application/json
// @Security TokenJWT
// @Param data body dto.CatalogRequest true "Catalog"
// @Success 200 {object} dto.CatalogRequest
// @Failure 400 {object} dto.Error Invalid JSON
// @Failure 500 {object} dto.Error Can't create catalog
// @Router /catalog [post]
func (cc *CatalogsController) CreateCatalog(c *gin.Context) {
	tracer := opentracing.GlobalTracer()
	controllerSpan := tracer.StartSpan("Controller:CreateCatalog")
	defer controllerSpan.Finish()
	ctx := context.Background()

	catalogDto := &dto.CatalogRequest{}
	if err := c.ShouldBindJSON(&catalogDto); err != nil {
		trace.OnError(cc.logger, controllerSpan, err)
		errs.ErrorHandler(c, errs.NewBadRequestError(errInvJSON+err.Error()))
		return
	}

	catalogResponse, err := cc.catalogsUC.CreateCatalog(ctx, catalogDto, controllerSpan)
	if err != nil {
		trace.OnError(cc.logger, controllerSpan, err)
		errs.ErrorHandler(c, errs.NewInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, catalogResponse)

}

// GetCatalogs godoc
// @Summary Get catalogs
// @Description Get JSON CatalogsRequest, return JSON GetCatalogsResponse
// @Tags Catalog
// @Produce  json
// @Content application/json
// @Security TokenJWT
// @Param data body dto.CatalogsRequest true "Catalogs"
// @Success 200 {object} dto.GetCatalogsResponse
// @Failure 400 {object} dto.Error Invalid JSON
// @Failure 500 {object} dto.Error Can't get catalogs
// @Router /catalog [get]
func (cc *CatalogsController) GetCatalogs(c *gin.Context) {
	tracer := opentracing.GlobalTracer()
	controllerSpan := tracer.StartSpan("Controller:GetCatalogs")
	defer controllerSpan.Finish()
	ctx := context.Background()

	catalogsDto := &dto.CatalogsRequest{}
	if err := c.ShouldBindJSON(&catalogsDto); err != nil {
		trace.OnError(cc.logger, controllerSpan, err)
		errs.ErrorHandler(c, errs.NewBadRequestError(errInvJSON+err.Error()))
		return
	}

	catalogsResponse, err := cc.catalogsUC.GetCatalogs(ctx, catalogsDto, controllerSpan)
	if err != nil {
		trace.OnError(cc.logger, controllerSpan, err)
		errs.ErrorHandler(c, errs.NewInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, catalogsResponse)
}

// GetCatalogCategories godoc
// @Summary Get categories
// @Description Return JSON GetCategoriesResponse
// @Tags Catalog
// @Produce  json
// @Content application/json
// @Security TokenJWT
// @Success 200 {object} dto.GetCategoriesResponse
// @Failure 500 {object} dto.Error Can't get categories
// @Router /categories [get]
func (cc *CatalogsController) GetCatalogCategories(c *gin.Context) {
	tracer := opentracing.GlobalTracer()
	controllerSpan := tracer.StartSpan("Controller:GetCatalogCategories")
	defer controllerSpan.Finish()
	ctx := context.Background()

	categoriesResponse, err := cc.catalogsUC.GetCatalogCategories(ctx, controllerSpan)
	if err != nil {
		trace.OnError(cc.logger, controllerSpan, err)
		errs.ErrorHandler(c, errs.NewInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, categoriesResponse)
}

// GetCatalogByID godoc
// @Summary Get catalog by ID
// @Description Get id in path, return JSON GetCatalogResponse
// @Tags Catalog
// @Produce  json
// @Content application/json
// @Security TokenJWT
// @Param id path int true "Catalog ID"
// @Success 200 {object} dto.GetCatalogResponse
// @Failure 400 {object} dto.Error Invalid ID
// @Failure 500 {object} dto.Error Can't get catalog
// @Router /catalog/:id [get]
func (cc *CatalogsController) GetCatalogByID(c *gin.Context) {
	tracer := opentracing.GlobalTracer()
	controllerSpan := tracer.StartSpan("Controller:GetCatalogByID")
	defer controllerSpan.Finish()
	ctx := context.Background()

	id, ok := c.Params.Get("id")
	if !ok || id == "" {
		trace.OnError(cc.logger, controllerSpan, errs.NewBadRequestError(errInvID))
		errs.ErrorHandler(c, errs.NewBadRequestError(errInvID))
		return
	}

	catalogResponse, err := cc.catalogsUC.GetCatalogByID(ctx, id, controllerSpan)
	if err != nil {
		trace.OnError(cc.logger, controllerSpan, err)
		errs.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, catalogResponse)
}

// UpdateCatalog godoc
// @Summary Update catalog
// @Description Get JSON CatalogRequest, return JSON UpdateCatalogResponse
// @Tags Catalog
// @Produce  json
// @Content application/json
// @Security TokenJWT
// @Param data body dto.CatalogRequest true "Catalog"
// @Success 200 {object} dto.UpdateCatalogResponse
// @Failure 400 {object} dto.Error Invalid JSON
// @Failure 500 {object} dto.Error Can't update catalog
// @Router /catalog [put]
func (cc *CatalogsController) UpdateCatalog(c *gin.Context) {
	tracer := opentracing.GlobalTracer()
	controllerSpan := tracer.StartSpan("Controller:UpdateCatalog")
	defer controllerSpan.Finish()
	ctx := context.Background()

	catalogDto := &dto.CatalogRequest{}
	if err := c.ShouldBindJSON(&catalogDto); err != nil {
		trace.OnError(cc.logger, controllerSpan, err)
		errs.ErrorHandler(c, errs.NewBadRequestError(errInvJSON+err.Error()))
		return
	}

	catalogResponse, err := cc.catalogsUC.UpdateCatalogByID(ctx, catalogDto, controllerSpan)
	if err != nil {
		trace.OnError(cc.logger, controllerSpan, err)
		errs.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, catalogResponse)
}

// DeleteCatalog godoc
// @Summary Delete catalog
// @Description Get id from path, return JSON DeleteCatalogResponse
// @Tags Catalog
// @Produce  json
// @Content application/json
// @Security TokenJWT
// @Param id path int true "Catalog ID"
// @Success 200 {object} dto.DeleteCatalogResponse
// @Failure 400 {object} dto.Error Invalid JSON
// @Failure 500 {object} dto.Error Can't delete catalog
// @Router /catalog/:id [delete]
func (cc *CatalogsController) DeleteCatalog(c *gin.Context) {
	tracer := opentracing.GlobalTracer()
	controllerSpan := tracer.StartSpan("Controller:DeleteCatalog")
	defer controllerSpan.Finish()
	ctx := context.Background()

	id, ok := c.Params.Get("id")
	if !ok || id == "" {
		trace.OnError(cc.logger, controllerSpan, errs.NewBadRequestError(errInvID))
		errs.ErrorHandler(c, errs.NewBadRequestError(errInvID))
		return
	}

	catalogResponse, err := cc.catalogsUC.DeleteCatalogByID(ctx, id, controllerSpan)
	if err != nil {
		trace.OnError(cc.logger, controllerSpan, err)
		errs.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, catalogResponse)
}
