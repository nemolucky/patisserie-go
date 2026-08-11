package auth_usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/nemolucky/patisserie-go/internal/domain"
)

func (u *UseCase) RegisterUser(ctx context.Context, input RegisterUserInput) (*domain.TokenPair, error) {

	ok, err := u.cache.IsVerified(ctx, input.Phone)
	if err != nil{
		return nil, err
	}

	if !ok{
		return nil, domain.ErrNotFound
	} 

	err = u.cache.DeleteVerified(ctx, input.Phone)
	if err != nil{
		return nil, err
	}

	names := strings.Fields(input.Fullname)
	if len(names) != 3{
		return nil, fmt.Errorf("Incorrect format: %s", names)
	}

	surname := names[0]
	name := names[1]
	patronymic := names[2]

	userID, err := u.repository.CreateUser(
			ctx,
			name,
			patronymic,
			surname,
			input.Phone,
			input.Telegram,
	)
	if err != nil{
		return nil, err
	}

	tokens, err := u.token.GeneratePairTokens(*userID)
	if err != nil{
		return nil, err
	}

	return tokens, nil
}