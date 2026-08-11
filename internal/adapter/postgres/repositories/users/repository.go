package users

import "github.com/nemolucky/patisserie-go/internal/adapter/postgres"

type UserRepository struct {
	*postgres.Repository
}

func New(base *postgres.Repository) *UserRepository {
	return &UserRepository{
		Repository: base,
	}
}