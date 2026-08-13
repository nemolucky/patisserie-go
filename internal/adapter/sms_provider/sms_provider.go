package smsprovider

import (
	"context"
	"fmt"
)

type SMSProvider struct{}

func (p *SMSProvider) SendSMS(ctx context.Context, phone string, code string) error{
	fmt.Println("Send SMS to ", phone, " with code ", code)
	return nil
}