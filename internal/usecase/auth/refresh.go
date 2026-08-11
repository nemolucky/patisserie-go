package auth_usecase

import (
	"context"

	"github.com/nemolucky/patisserie-go/internal/domain"
)

func (u *UseCase) Refresh(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	claims, err := u.token.Parse(refreshToken, domain.RefreshTokenType)
	if err != nil {
		return nil, err
	}

	return u.token.GeneratePairTokens(claims.UserID)
}