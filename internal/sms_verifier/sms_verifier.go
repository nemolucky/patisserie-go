package sms_verifier

import "context"

type Cache interface{
	GetSMSCode(ctx context.Context, phone string) (string, error)
	SetSMSCode(ctx context.Context, phone string, code string) (error)
	DeleteSMSCode(ctx context.Context, phone string) (error)
}

type SMSProvider interface{
	SendSMS(ctx context.Context, phone string, code string) error
}

type SMSVerifier struct{
	sms SMSProvider
	cache Cache
}

func New(
	sms SMSProvider,
	cache Cache,
) *SMSVerifier {
	return &SMSVerifier{
		sms:   sms,
		cache: cache,
	}
}

func (s *SMSVerifier) SendCode(ctx context.Context, phone string) error{
		
	code := GenerateSMSCode()

	err := s.cache.SetSMSCode(ctx, phone, code)
	if err != nil{
		return err
	}

	err = s.sms.SendSMS(ctx, phone, code)
	if err != nil{
		return err
	}

	return nil
}

func (s *SMSVerifier) VerifyCode(ctx context.Context, phone string, code string) (bool, error){
	cachedCode, err := s.cache.GetSMSCode(ctx, phone)
	if err != nil {
		return false, err
	}

	if cachedCode != code{
		return false, nil
	}

	err = s.cache.DeleteSMSCode(ctx, phone)
	if err != nil{
		return false, err
	}

	return true, nil
}