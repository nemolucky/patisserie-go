package category_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

type Repository interface {
	ListCategory(ctx context.Context) ([]domain.Category, error)
	CreateCategory(ctx context.Context, title string, parent uuid.UUID) (uuid.UUID, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, title string) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
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


func (u *UseCase) Create(ctx context.Context, input CreateInput) (uuid.UUID, error) {
	return u.repository.CreateCategory(ctx, input.Title, input.Parent)
}

func (u *UseCase) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*domain.Category, error) {
	return u.repository.UpdateCategory(ctx, id, input.Title)
}

func (u *UseCase) Delete(ctx context.Context, input DeleteInput) error {
	return u.repository.DeleteCategory(ctx, input.ID)
}

func (u *UseCase) List(ctx context.Context) ([]domain.Category, error) {
	return u.repository.ListCategory(ctx)
}

func (u *UseCase) Tree(ctx context.Context) ([]*CategoryNode, error) {
	
	categories, err := u.repository.ListCategory(ctx)
	if err != nil {
		return nil, err
	}

	return BuildTree(categories), nil
}