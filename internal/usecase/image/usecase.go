package image_usecase

import (
	"context"

	"github.com/nemolucky/patisserie-go/internal/domain"
)

type Storage interface {
	Upload(ctx context.Context, image *domain.UploadImage) (*domain.UploadedImage, error)
}

type UseCase struct {
	Storage Storage
}

func New(
	storage Storage,
) *UseCase {
	return &UseCase{
		Storage: storage,
	}
}

func (u *UseCase) Upload(ctx context.Context, image *domain.UploadImage) (*domain.UploadedImage, error) {
	return u.Storage.Upload(ctx, image)
}