package image

import "github.com/nemolucky/patisserie-go/internal/adapter/postgres"


type ImageRepository struct {
	*postgres.Repository
}

func New(base *postgres.Repository) *ImageRepository {
	return &ImageRepository{
		Repository: base,
	}
}