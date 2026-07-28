package products

import (
	"context"
	"encoding/json"

	_ "embed"

	"github.com/google/uuid"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

//go:embed query/list_catalog.sql
var ListCatalogProductsQuery string

//go:embed query/create.sql
var CreateProductsQuery string

//go:embed query/connect_images.sql
var ConnectProductImagesQuery string

//go:embed query/delete_attached_images.sql
var DeleteAttachedImagesQuery string

//go:embed query/connect_categories.sql
var ConnectProductCategoriesQuery string

//go:embed query/delete_attached_categories.sql
var DeleteAttachedCategoriesQuery string

//go:embed query/create_nutrition.sql
var ConnectProductNutritionQuery string

//go:embed query/update.sql
var UpdateProductQuery string

//go:embed query/update_nutrition.sql
var UpdateNutritionQuery string

//go:embed query/get_by_id.sql
var GetProductByIDQuery string

//go:embed query/toggle_active.sql
var ToggleActiveByIDQuery string


func (r *ProductRepository) ListCatalogProducts(ctx context.Context,  search string, limit int, offset int) ([]*domain.CatalogProduct, error) {

	products := make([]*domain.CatalogProduct, 0, limit)

	rows, err := r.Query(ctx, ListCatalogProductsQuery, search, limit, offset)
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

func (r *ProductRepository) CreateProduct(ctx context.Context,  title, description string, price int64) (*uuid.UUID, error) {

	var productID uuid.UUID

	err := r.QueryRow(ctx, CreateProductsQuery, title, description, price).Scan(&productID)
	if err != nil{
		return nil, err
	}

	return &productID, nil
}

func (r *ProductRepository) ConnectProductImages(ctx context.Context,  productID uuid.UUID, images []uuid.UUID) error {

	orders := make([]int, len(images))

	for i := range images {
		orders[i] = i
	}

	_, err := r.Exec(ctx, ConnectProductImagesQuery, productID, images, orders)

	return err
}

func (r *ProductRepository) DeleteAttachedImages(ctx context.Context,  id uuid.UUID) error {

	_, err := r.Exec(ctx, DeleteAttachedImagesQuery, id)
	
	return err
}

func (r *ProductRepository) ConnectProductCategories(ctx context.Context,  productID uuid.UUID, categories []uuid.UUID) error {

	_, err := r.Exec(ctx, ConnectProductCategoriesQuery, productID, categories)

	return err
}


func (r *ProductRepository) DeleteAttachedCategories(ctx context.Context,  id uuid.UUID) error {

	_, err := r.Exec(ctx, DeleteAttachedCategoriesQuery, id)
	
	return err
}

func (r *ProductRepository) CreateNutrition(ctx context.Context,  productID uuid.UUID, nutrition domain.NutritionFacts) error {

	_, err := r.Exec(ctx, ConnectProductNutritionQuery, productID, nutrition.Calories, nutrition.Protein, nutrition.Fat, nutrition.Carbohydrates, nutrition.Weight)
	
	return err
}

func (r *ProductRepository) GetProduct(ctx context.Context,  id uuid.UUID) (*domain.Product, error) {

	var product domain.Product

	var (
		categories json.RawMessage
		images json.RawMessage
		facts json.RawMessage
	)

	err := r.QueryRow(ctx, GetProductByIDQuery, id).
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

func (r *ProductRepository) ToggleActiveProduct(ctx context.Context,  id uuid.UUID) (bool, error) {

	var isActive bool
	err := r.QueryRow(ctx, ToggleActiveByIDQuery, id).Scan(&isActive)
	
	return isActive, err
}

func (r *ProductRepository) UpdateProduct(ctx context.Context,  id uuid.UUID, title, description string, price int64) error {

	_ , err := r.Exec(ctx, UpdateProductQuery, id, title, description, price)

	return err
}

func (r *ProductRepository) UpdateNutrition(ctx context.Context, id uuid.UUID, facts domain.NutritionFacts) error {

	_ , err := r.Exec(ctx, UpdateNutritionQuery, id, facts.Calories, facts.Protein, facts.Fat, facts.Carbohydrates, facts.Weight)

	return err
}
