package product_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

func (u *UseCase) UpdateProduct(ctx context.Context, id uuid.UUID, input UpdateProductInput) error {
	return u.repository.UpdateProduct(ctx, id, input.Title, input.Description, input.Price)
}

func (u *UseCase) UpdateNutrition(ctx context.Context, productID uuid.UUID, facts domain.NutritionFacts) error {
	return u.repository.UpdateNutrition(ctx, productID, facts)
}

func (u *UseCase) UpdateAttachedCategories(ctx context.Context, productID uuid.UUID, categories []uuid.UUID) error {
	
	err := u.tm.Do(ctx, func(ctx context.Context) error {
		if err := u.repository.DeleteAttachedCategories(ctx, productID); err != nil {
				return err
		}

		if err := u.repository.ConnectProductCategories(ctx, productID, categories); err != nil {
				return err
		}

		return nil
	})

	return err
}

func (u *UseCase) UpdateAttachedImages(ctx context.Context, productID uuid.UUID, images []uuid.UUID) error {

	err := u.tm.Do(ctx, func(ctx context.Context) error {
		if err := u.repository.DeleteAttachedImages(ctx, productID); err != nil {
				return err
		}

		if err := u.repository.ConnectProductImages(ctx, productID, images); err != nil {
				return err
		}

		return nil
	})

	return err
}