package product_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

type Repository interface{
	BeginTx(ctx context.Context) (pgx.Tx, error)
	CommitTx(ctx context.Context, tx pgx.Tx) error
	ListCatalogProducts(ctx context.Context, limit int, offset int) ([]*domain.CatalogProduct, error)
	CreateProduct(ctx context.Context, title, description string, price int64) (*uuid.UUID, error)
	CreateProductImages(ctx context.Context, productID uuid.UUID, images []ProductImageInput) error
	CreateProductCategories(ctx context.Context, productID uuid.UUID, categories []uuid.UUID) error
	CreateNutrition(ctx context.Context, productID uuid.UUID, nutrition domain.NutritionFacts) error
	GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error)
}

type UseCase struct {
	repository Repository
}

func New(
	repository Repository,
) *UseCase {
	return &UseCase{
		repository: repository,
	}
}

func (u *UseCase) ListCatalogProducts(ctx context.Context, input ListCatalogProductsInput) (*ListCatalogProductsOutput, error) {
	
	limit := input.Limit
	offset := input.Offset

	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	products, err := u.repository.ListCatalogProducts(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return &ListCatalogProductsOutput{
		Products: products,
		Limit:    limit,
		Offset:   offset,
	}, nil
}

func (u *UseCase) Create(ctx context.Context, input CreateProductInput) (*uuid.UUID, error) {

	tx, err := u.repository.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	productID, err := u.repository.CreateProduct(ctx, input.Title, input.Description, input.Price)
	if err != nil {
		return nil, err
	}

	err = u.repository.CreateProductImages(ctx, *productID, input.Images)
	if err != nil {
		return nil, err
	}

	err = u.repository.CreateProductCategories(ctx, *productID, input.Categories)
	if err != nil {
		return nil, err
	}

	err = u.repository.CommitTx(ctx, tx)

	if err != nil {
		return nil, err
	}

	return productID, nil
}

func (u *UseCase) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return u.repository.GetProduct(ctx, id)
}