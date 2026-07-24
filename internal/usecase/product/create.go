package product_usecase

import (
	"context"

	"github.com/google/uuid"
)

func (u *UseCase) Create(ctx context.Context, input CreateProductInput) (*uuid.UUID, error) {

	tx, err := u.repository.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	productID, err := u.repository.CreateProduct(ctx, input.Title, input.Description, input.Price)
	if err != nil {
		return nil, err
	}

	err = u.repository.ConnectProductImages(ctx, *productID, input.Images)
	if err != nil {
		return nil, err
	}

	err = u.repository.ConnectProductCategories(ctx, *productID, input.Categories)
	if err != nil {
		return nil, err
	}

	err = u.repository.CreateNutrition(ctx, *productID, *input.Facts)
	if err != nil {
		return nil, err
	}

	err = u.repository.CommitTx(ctx, tx)

	if err != nil {
		return nil, err
	}

	return productID, nil
}