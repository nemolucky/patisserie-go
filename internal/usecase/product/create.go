package product_usecase

import (
	"context"

	"github.com/google/uuid"
)

func (u *UseCase) Create(ctx context.Context, input CreateProductInput) (*uuid.UUID, error) {

	var productID *uuid.UUID

	err := u.tm.Do(ctx, func(ctx context.Context) error {
		var err error

		productID, err = u.repository.CreateProduct(ctx, input.Title, input.Description, input.Price)
		if err != nil {
			return err
		}

		err = u.repository.ConnectProductImages(ctx, *productID, input.Images)
		if err != nil {
			return err
		}

		err = u.repository.ConnectProductCategories(ctx, *productID, input.Categories)
		if err != nil {
			return err
		}

		err = u.repository.CreateNutrition(ctx, *productID, input.Facts)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return productID, nil
}