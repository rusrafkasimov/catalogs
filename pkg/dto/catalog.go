package dto

type CatalogRequest struct {
	ID       string `json:"id"`
	Active   bool   `json:"active"`
	Category string `json:"category"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Value    string `json:"value"`
} // @Name CatalogRequest

type CatalogsRequest struct {
	Category string `query:"category" json:"category" validate:"required"`
	Query    string `query:"query" json:"query"`
	Sorted   bool   `query:"sorted" json:"sorted" default:"true"`
}
