package auth_usecase

import "github.com/nemolucky/patisserie-go/internal/domain"

type SendSMSCodeInput struct {
	Phone    	string `json:"phone"`
}

type VerifySMSCodeInput struct {
	Phone    	string `json:"phone"`
	Code    	string `json:"code"`
}

type VerifySMSOutput struct {
	RegistrationRequired bool 		`json:"registration_required"`
	TokenPair domain.TokenPair `json:"token_pair"`
}

type RegisterUserInput struct {
	Fullname  string `json:"full_name"`
	Phone    	string `json:"phone"`
	Telegram *string `json:"telegram"`
}
