package auth_usecase

import (
	"context"
	"fmt"
)

func (u *UseCase) VerifySMSCode(ctx context.Context, input VerifySMSCodeInput)(*VerifySMSOutput, error){

	phone := input.Phone
	code := input.Code

	ok, err := u.sms.VerifyCode(ctx, phone, code)
	if err != nil{
		return nil, err
	}

	if !ok{
		return nil, fmt.Errorf("invalid code for phone: %s", phone)
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