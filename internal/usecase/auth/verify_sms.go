package auth_usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/nemolucky/patisserie-go/internal/domain"
)

func (u *UseCase) VerifySMSCode(ctx context.Context, input VerifySMSCodeInput)(*VerifySMSOutput, error){

	phone := input.Phone
	code := input.Code

	cachedCode, err := u.cache.GetSMSCode(ctx, phone)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrExpired
		}
		return nil, err
	}

	if cachedCode != code{
		return nil, fmt.Errorf("Incorrect code: %s", code)
	}

	err = u.cache.DeleteSMSCode(ctx, phone)
	if err != nil{
		return nil, err
	}


	user, err := u.repository.GetUserByPhone(ctx, phone)
	if err != nil{
		return nil, err
	}

	if user == nil{

		err = u.cache.SetVerified(ctx, phone)
		if err != nil{
			return nil, err
		}

		return &VerifySMSOutput{
			RegistrationRequired: true,
		}, err
	}else{

		tokens, err := u.token.GeneratePairTokens(user.ID)
		if err != nil{
			return nil, err
		}

		return &VerifySMSOutput{
			TokenPair: *tokens,
			RegistrationRequired: false,
		}, nil
	}
}