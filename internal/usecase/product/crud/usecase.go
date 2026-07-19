package product_usecase

import (
	"context"

	"github.com/nemolucky/patisserie-go/internal/domain"
)

type Postgres interface{
	ListCatalogProducts(ctx context.Context, limit int, offset int) ([]*domain.CatalogProduct, error)
}

type UseCase struct {
	Postgres Postgres
}

func New(
	postgres Postgres,
) *UseCase {
	return &UseCase{
		Postgres: postgres,
	}
}

func (u *UseCase) ListCatalogProducts(ctx context.Context, input ListCatalogProductsInput) ([]*domain.CatalogProduct, error) {
	return u.Postgres.ListCatalogProducts(ctx, input.Limit, input.Offset)
}