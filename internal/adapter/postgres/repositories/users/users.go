package users

import (
	"context"

	_ "embed"

	"github.com/google/uuid"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

//go:embed query/create.sql
var CreateUserQuery string

//go:embed query/get_by_phone.sql
var GetUserByPhoneQuery string

func (r *UserRepository) CreateUser(ctx context.Context, firstName, MiddleName, LastName, phone string, telegram *string) (*uuid.UUID, error) {
	
	var userID *uuid.UUID

	err := r.QueryRow(ctx, CreateUserQuery, firstName, MiddleName, LastName, phone, telegram).Scan(&userID)
	if err != nil{
		return nil, err
	}
	
	return userID, nil
}

func (r *UserRepository) GetUserByPhone(ctx context.Context, phone string) (*domain.User, error) {
	
	var user domain.User

	err := r.QueryRow(ctx, GetUserByPhoneQuery, phone).
		Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.Phone, &user.Telegram)
	if err != nil{
		return nil, err
	}
	
	return &user, nil
}