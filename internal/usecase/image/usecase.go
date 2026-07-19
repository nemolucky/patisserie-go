package image_usecase

import (
	"context"

	"github.com/nemolucky/patisserie-go/internal/domain"
)

type Storage interface {
	Upload(ctx context.Context, image UploadImage) (*domain.UploadedImage, error)
	Delete(ctx context.Context, key string) error
}

type Repository interface {
	CreateImage(ctx context.Context, key string) error
}

type UseCase struct {
	storage Storage
	repository Repository
}

func New(
	storage Storage,
	repository Repository,
) *UseCase {
	return &UseCase{
		storage: storage,
		repository: repository,
	}
}

func (u *UseCase) Upload(ctx context.Context, images []UploadImage) ([]*domain.UploadedImage, error) {

	uploadedImages := make([]*domain.UploadedImage, 0, len(images))

	for _, image := range images {
		uploadedImage, err := u.storage.Upload(ctx, image)
		if err != nil {
			return nil, err
		}
		
		err = u.repository.CreateImage(ctx, uploadedImage.Key)
		if err != nil {
			_ = u.storage.Delete(ctx, uploadedImage.Key)
			return nil, err
		}

		uploadedImages = append(uploadedImages, uploadedImage)
	}
	
	return uploadedImages, nil
}