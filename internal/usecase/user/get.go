package user_usecase

import (
	"context"

	"github.com/nemolucky/patisserie-go/internal/ctxutil"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

func (u *UseCase) GetUserByPhone(ctx context.Context, phone string) (*domain.User, error) {
	return u.repository.GetUserByPhone(ctx, phone)
}

func (u *UseCase) GetMe(ctx context.Context) (*domain.User, error) {

	userID, ok := ctxutil.UserIDFromContext(ctx)
	if !ok  {
		return nil, domain.ErrNotFound
	}

	return u.repository.GetUserByID(ctx, userID)
}
