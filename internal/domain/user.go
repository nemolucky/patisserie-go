package domain

import "github.com/google/uuid"

type User struct {
	ID       				uuid.UUID `json:"id"`
	FirstName       string `json:"first_name"`
	MiddleName      string `json:"middle_name"`
	LastName        string `json:"last_name"`
	Phone       		string `json:"phone"`
	Telegram       	*string `json:"telegram"`
}