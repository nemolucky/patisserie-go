package auth_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

type Repository interface{
	CreateUser(ctx context.Context, firstName, MiddleName, LastName, phone string, telegram *string) (*uuid.UUID, error)
	GetUserByPhone(ctx context.Context, phone string) (*domain.User, error)
}

type Cache interface{
	GetSMSCode(ctx context.Context, phone string) (string, error)
	SetSMSCode(ctx context.Context, phone string, code string) (error)
	DeleteSMSCode(ctx context.Context, phone string) (error)
	SetVerified(ctx context.Context, phone string) (error)
	IsVerified(ctx context.Context, phone string) (bool, error)
	DeleteVerified(ctx context.Context, phone string) error
}

type SMSProvider interface{
	SendSMS(ctx context.Context, phone string, code string) (bool, error)
}

type TokenProvider interface{
	GeneratePairTokens(userID uuid.UUID) (*domain.TokenPair, error)
	Parse(tokenString string, tokenType domain.TokenType) (*domain.Claims, error)
}

type UseCase struct{
	repository Repository
	cache Cache
	sms SMSProvider
	token TokenProvider
}

func New(
	repository Repository,
	cache Cache,
	sms SMSProvider,
	token TokenProvider,
) *UseCase {
	return &UseCase{
		repository: repository,
		cache: cache,
		sms: sms,
		token: token,
	}
}