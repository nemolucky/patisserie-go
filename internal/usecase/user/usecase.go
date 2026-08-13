package user_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

type Repository interface{
	CreateUser(ctx context.Context, firstName, MiddleName, LastName, phone string, telegram *string) (*uuid.UUID, error)
	GetUserByPhone(ctx context.Context, phone string) (*domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	UpdatePhone(ctx context.Context, id uuid.UUID, phone string) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type Cache interface{
	GetSMSCode(ctx context.Context, phone string) (string, error)
	SetSMSCode(ctx context.Context, phone string, code string) (error)
	DeleteSMSCode(ctx context.Context, phone string) (error)
}

type SMSVerifier interface{
	SendCode(ctx context.Context, phone string) error
	VerifyCode(ctx context.Context, phone string, code string) (bool, error)
}

type UseCase struct{
	repository Repository
	cache Cache
	sms SMSVerifier
}

func New(
	repository Repository,
	cache Cache,
	sms SMSVerifier,
) *UseCase {
	return &UseCase{
		repository: repository,
		cache: cache,
		sms: sms,
	}
}