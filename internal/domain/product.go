package domain

import "github.com/google/uuid"

type NutritionFacts struct {
	Calories      int     `json:"calories"`
	Protein       float32 `json:"protein"`
	Fat           float32 `json:"fat"`
	Carbohydrates float32 `json:"carbohydrates"`
	Weight        float32 `json:"weight"`
}

type ProductImage struct {
	ID        uuid.UUID `json:"id"`
	URL       string    `json:"url"`
	SortOrder int       `json:"sort_order"`
}

type Product struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	IsActive    bool      `json:"is_active"`

	Facts      *NutritionFacts `json:"facts,omitempty"`
	Images     []ProductImage  `json:"images"`
	Categories []Category      `json:"categories"`
}

type CatalogProduct struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Price       int64     `json:"price"`
	Image     ProductImage  `json:"image"`
	Categories []Category      `json:"categories"`
}

// type UpdateProductInput struct {
// 	ID          uuid.UUID
// 	Title       *string
// 	Description *string
// 	Price       *int64
// }
