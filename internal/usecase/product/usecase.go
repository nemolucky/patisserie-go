package product_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

type Repository interface{

	ListCatalogProducts(ctx context.Context, search string, limit int, offset int) ([]*domain.CatalogProduct, error)

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

type TransactionManager interface{
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type UseCase struct {
	repository Repository
	tm TransactionManager
}

func New(
	repository Repository,
	tm TransactionManager,
) *UseCase {
	return &UseCase{
		repository: repository,
		tm: tm,
	}
}

func (u *UseCase) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return u.repository.GetProduct(ctx, id)
}

func (u *UseCase) ToggleActiveProduct(ctx context.Context, id uuid.UUID) (bool, error) {
	return u.repository.ToggleActiveProduct(ctx, id)
}