package product_usecase

import "github.com/nemolucky/patisserie-go/internal/domain"

type ListCatalogProductsInput struct{
	Limit int `json:"limit"`
	Offset int `json:"offset"`
}

type ListCatalogProductsOutput struct{
	Products []*domain.CatalogProduct `json:"products"`
	Limit int `json:"limit"`
	Offset int `json:"offset"`
}