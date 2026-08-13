package user_usecase

import (
	"context"
	"fmt"

	"github.com/nemolucky/patisserie-go/internal/ctxutil"
)

func (u *UseCase) SendSMSNewPhone(ctx context.Context, input SendSMSNewPhoneInput) error{
	return u.sms.SendCode(ctx, input.Phone)
}

func (u *UseCase) UpdatePhone(ctx context.Context, input UpdatePhoneInput) error{
	isValid, err := u.sms.VerifyCode(ctx, input.Phone, input.Code)
	if err != nil {
		return err
	}
	if !isValid {
		return fmt.Errorf("invalid code for phone: %s", input.Phone)
	}

	userID, ok := ctxutil.UserIDFromContext(ctx)
	if !ok  {
		return fmt.Errorf("user not found")
	}

	return u.repository.UpdatePhone(ctx, userID, input.Phone)	
}
