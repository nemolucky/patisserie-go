package user_usecase

type SendSMSNewPhoneInput struct {
	Phone string `json:"phone"`
}

type UpdatePhoneInput struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}