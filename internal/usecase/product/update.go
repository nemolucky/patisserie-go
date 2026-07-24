package product_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

func (u *UseCase) UpdateProduct(ctx context.Context, id uuid.UUID, input UpdateProductInput) error {
	return u.repository.UpdateProduct(ctx, id, input.Title, input.Description, input.Price)
}

func (u UseCase) UpdateNutrition(ctx context.Context, productID uuid.UUID, facts domain.NutritionFacts) error {
	return u.repository.UpdateNutrition(ctx, productID, facts)
}

func (u *UseCase) UpdateAttachedCategories(ctx context.Context, productID uuid.UUID, categories []uuid.UUID) error {
	
	tx, err := u.repository.BeginTx(ctx)
	if err != nil {
		return err
	}

	err = u.repository.DeleteAttachedCategories(ctx, productID)
	if err != nil {
		return err
	}

	err = u.repository.ConnectProductCategories(ctx, productID, categories)
	if err != nil {
		return err
	}

	return u.repository.CommitTx(ctx, tx)
}

func (u *UseCase) UpdateAttachedImages(ctx context.Context, productID uuid.UUID, images []uuid.UUID) error {
	
	tx, err := u.repository.BeginTx(ctx)
	if err != nil {
		return err
	}

	err = u.repository.DeleteAttachedImages(ctx, productID)
	if err != nil {
		return err
	}

	err = u.repository.ConnectProductImages(ctx, productID, images)
	if err != nil {
		return err
	}

	return u.repository.CommitTx(ctx, tx)
}