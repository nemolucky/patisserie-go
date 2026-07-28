package products

import "github.com/nemolucky/patisserie-go/internal/adapter/postgres"


type ProductRepository struct {
	*postgres.Repository
}

func New(base *postgres.Repository) *ProductRepository {
	return &ProductRepository{
		Repository: base,
	}
}