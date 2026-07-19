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

func (u *UseCase) Upload(ctx context.Context, images []*domain.UploadImage) ([]*domain.UploadedImage, error) {

	uploadedImages := make([]*domain.UploadedImage, 0, len(images))

	for _, image := range images {
		uploadedImage, err := u.Storage.Upload(ctx, image)
		if err != nil {
			return nil, err
		}
		uploadedImages = append(uploadedImages, uploadedImage)
	}
	
	return uploadedImages, nil
}