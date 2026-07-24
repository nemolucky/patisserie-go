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
	ConnectProductImages(ctx context.Context, productID uuid.UUID, images []uuid.UUID) error
	ConnectProductCategories(ctx context.Context, productID uuid.UUID, categories []uuid.UUID) error
	CreateNutrition(ctx context.Context, productID uuid.UUID, nutrition domain.NutritionFacts) error

	UpdateProduct(ctx context.Context, id uuid.UUID, title, description string, price int64) error
	UpdateNutrition(ctx context.Context, productID uuid.UUID, facts domain.NutritionFacts) error

	DeleteAttachedCategories(ctx context.Context, productID uuid.UUID) error
	DeleteAttachedImages(ctx context.Context, productID uuid.UUID) error

	GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	ToggleActiveProduct(ctx context.Context, id uuid.UUID) (bool, error)
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

func (u *UseCase) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return u.repository.GetProduct(ctx, id)
}

func (u *UseCase) ToggleActiveProduct(ctx context.Context, id uuid.UUID) (bool, error) {
	return u.repository.ToggleActiveProduct(ctx, id)
}