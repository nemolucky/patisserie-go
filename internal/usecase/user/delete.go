package user_usecase

import (
	"context"

	"github.com/nemolucky/patisserie-go/internal/ctxutil"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

func (u *UseCase) DeleteAccount(ctx context.Context) error {

	userID, ok := ctxutil.UserIDFromContext(ctx)
	if !ok  {
		return domain.ErrNotFound
	}

	return u.repository.DeleteUser(ctx, userID)
}