package auth_usecase

import "context"

func (u *UseCase) SendSMSCode(ctx context.Context, input SendSMSCodeInput)(error){

	phone := input.Phone

	err := u.sms.SendCode(ctx, phone)
	if err != nil{
		return err
	}

	return nil
}