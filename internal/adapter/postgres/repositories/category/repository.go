package category

import "github.com/nemolucky/patisserie-go/internal/adapter/postgres"


type CategoryRepository struct {
	*postgres.Repository
}

func New(base *postgres.Repository) *CategoryRepository {
	return &CategoryRepository{
		Repository: base,
	}
}