package user_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

type Repository interface{
	CreateUser(ctx context.Context, firstName, MiddleName, LastName, phone string, telegram *string) (*uuid.UUID, error)
	GetUserByPhone(ctx context.Context, phone string) (*domain.User, error)
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



func (u *UseCase) GetUserByPhone(ctx context.Context, phone string) (*domain.User, error) {
	return u.repository.GetUserByPhone(ctx, phone)
}