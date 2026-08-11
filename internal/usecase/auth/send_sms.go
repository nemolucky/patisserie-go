package auth_usecase

import "context"

func (u *UseCase) SendSMSCode(ctx context.Context, input SendSMSCodeInput)(error){

	phone := input.Phone

	code := GenerateSMSCode()

	err := u.cache.SetSMSCode(ctx, phone, code)
	if err != nil{
		return err
	} 

	if ok, err := u.sms.SendSMS(ctx, phone, code); !ok{
		return err
	}

	return nil
}