package product_usecase

import (
	"github.com/google/uuid"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

type ListCatalogProductsInput struct{
	Limit int `json:"limit"`
	Offset int `json:"offset"`
}

type ListCatalogProductsOutput struct{
	Products []*domain.CatalogProduct `json:"products"`
	Limit int `json:"limit"`
	Offset int `json:"offset"`
}

type ProductImageInput struct {
	ID       uuid.UUID  `json:"id"`
	SortOrder int       `json:"sort_order"`
}

type CreateProductInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int64 `json:"price"`

	Images     []ProductImageInput `json:"images"`
	Categories []uuid.UUID `json:"categories"`
	Facts      *domain.NutritionFacts `json:"facts,omitempty"`
}