package image_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

type Storage interface {
	Upload(ctx context.Context, image UploadImage) (*domain.UploadedImage, error)
	Delete(ctx context.Context, key string) error
}

type Repository interface {
	CreateImage(ctx context.Context, key string) (*uuid.UUID, error)
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

func (u *UseCase) Upload(ctx context.Context, images []UploadImage) ([]*uuid.UUID, error) {

	uploadedImages := make([]*uuid.UUID, 0, len(images))

	for _, image := range images {
		uploadedImage, err := u.storage.Upload(ctx, image)
		if err != nil {
			return nil, err
		}
		
		imageID, err := u.repository.CreateImage(ctx, uploadedImage.Key)
		if err != nil {
			_ = u.storage.Delete(ctx, uploadedImage.Key)
			return nil, err
		}

		uploadedImages = append(uploadedImages, imageID)
	}
	
	return uploadedImages, nil
}