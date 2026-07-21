package postgres

import (
	"context"
	"encoding/json"

	_ "embed"

	"github.com/google/uuid"
	"github.com/nemolucky/patisserie-go/internal/domain"
	product_usecase "github.com/nemolucky/patisserie-go/internal/usecase/product/crud"
)

//go:embed queries/products/list_catalog.sql
var ListCatalogProductsQuery string

//go:embed queries/products/create.sql
var CreateProductsQuery string

//go:embed queries/products/connect_images.sql
var ConnectProductImagesQuery string

//go:embed queries/products/connect_categories.sql
var ConnectProductCategoriesQuery string

//go:embed queries/products/create_nutrition.sql
var ConnectProductNutritionQuery string

//go:embed queries/products/get_by_id.sql
var GetProductByIDQuery string


func (p *Postgres) ListCatalogProducts(ctx context.Context, limit int, offset int) ([]*domain.CatalogProduct, error) {

	products := make([]*domain.CatalogProduct, 0, limit)

	rows, err := p.Pool.Query(ctx, ListCatalogProductsQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var product domain.CatalogProduct

		if err := rows.Scan(&product.ID, &product.Title, &product.Price, &product.Image.Key, &product.Image.SortOrder, &product.Categories); err != nil {
			return nil, err
		}

		products = append(products, &product)
	}
	if err := rows.Err(); err != nil {
    return nil, err
	}


	return products, nil
}

func (p *Postgres) CreateProduct(ctx context.Context, title, description string, price int64) (*uuid.UUID, error) {

	var productID uuid.UUID

	err := p.Pool.QueryRow(ctx, CreateProductsQuery, title, description, price).Scan(&productID)
	if err != nil{
		return nil, err
	}

	return &productID, nil
}

func (p *Postgres) CreateProductImages(ctx context.Context, productID uuid.UUID, images []product_usecase.ProductImageInput) error {

	for _, image := range images {
		_, err := p.Pool.Exec(ctx, ConnectProductImagesQuery, productID, image.ID, image.SortOrder)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Postgres) CreateProductCategories(ctx context.Context, productID uuid.UUID, categories []uuid.UUID) error {

	for _, category := range categories {
		_, err := p.Pool.Exec(ctx, ConnectProductCategoriesQuery, productID, category)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Postgres) CreateNutrition(ctx context.Context, productID uuid.UUID, nutrition domain.NutritionFacts) error {

	_, err := p.Pool.Exec(ctx, ConnectProductNutritionQuery, productID, nutrition.Calories, nutrition.Protein, nutrition.Fat, nutrition.Carbohydrates, nutrition.Weight)
	
	return err
}

func (p *Postgres) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error) {

	var product domain.Product

	var (
		categories json.RawMessage
		images json.RawMessage
		facts json.RawMessage
	)

	err := p.Pool.QueryRow(ctx, GetProductByIDQuery, id).
								Scan(&product.ID, 
										&product.Title, 
										&product.Description, 
										&product.Price, 
										&product.IsActive, 
										&categories, 
										&images, 
										&facts)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(categories, &product.Categories); err != nil {
			return nil, err
	}

	if err := json.Unmarshal(images, &product.Images); err != nil {
			return nil, err
	}

	if len(facts) > 0 && string(facts) != "null" {
			if err := json.Unmarshal(facts, &product.Facts); err != nil {
					return nil, err
			}
	}

	return &product, nil
}