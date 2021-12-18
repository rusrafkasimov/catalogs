package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type CatalogResponse struct {
	ID       primitive.ObjectID `json:"id"`
	Active   bool   `json:"active"`
	Category string `json:"category"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Value    string `json:"value"`
} // @Name CatalogResponse

type CreateCatalogResponse struct {
	Payload struct {
		CatalogResponse `mapstructure:",squash"`
	} `json:"payload"`
	Meta ResponseMeta `json:"meta"`
}// @Name CreateCatalogResponse

type GetCatalogResponse struct {
	Payload struct {
		CatalogResponse `mapstructure:",squash"`
	} `json:"payload"`
	Meta ResponseMeta `json:"meta"`
}// @Name GetCatalogResponse

type GetCatalogsResponse struct {
	Payload []CatalogResponse `json:"payload" mapstructure:",squash"`
	Meta    ResponseMetaList  `json:"meta"`
}// @Name GetCatalogsResponse

type GetCategoriesResponse struct {
	Payload []string         `json:"payload" mapstructure:",squash"`
	Meta    ResponseMetaList `json:"meta"`
}// @Name GetCategoriesResponse

type UpdateCatalogResponse struct {
	Payload struct {
		CatalogResponse `mapstructure:",squash"`
	} `json:"payload"`
	Meta ResponseMeta `json:"meta"`
}// @Name UpdateCatalogResponse

type DeleteCatalogResponse struct {
	Payload struct {
		Active bool `json:"active"`
	} `json:"payload"`
	Meta ResponseMeta `json:"meta"`
}// @Name DeleteCatalogResponse