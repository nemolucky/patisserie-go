package smsprovider

import (
	"context"
	"fmt"
)

type SMSProvider struct{}

func (p *SMSProvider) SendSMS(ctx context.Context, phone string, code string) (bool, error){
	fmt.Println("Send SMS to ", phone, " with code ", code)
	return true, nil
}