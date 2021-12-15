package controllers

import (
	"github.com/afiskon/promtail-client/promtail"
	"github.com/rusrafkasimov/catalogs/pkg/usecases"
)

type CatalogsController struct {
	logger promtail.Client
	catalogsUC usecases.CatalogsUseCase
}

func NewCatalogsController(uc usecases.CatalogsUseCase, logger promtail.Client) *CatalogsController {
	return &CatalogsController{
		logger: logger,
		catalogsUC: uc,
	}
}