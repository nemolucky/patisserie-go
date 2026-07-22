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
	DeleteImages(ctx context.Context, imageIDs []uuid.UUID) error
	FindUnusedImages(ctx context.Context) ([]*domain.UnusedImage, error)
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

func (u *UseCase) DeleteUnusedImages(ctx context.Context) ([]*domain.UnusedImage, error) {

	unusedImages, err := u.repository.FindUnusedImages(ctx)
	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, 0, len(unusedImages))

	for _, image := range unusedImages {

		err = u.storage.Delete(ctx, image.Key)
		if err != nil {
			return nil, err
		}

		ids = append(ids, image.ID)
	}

	if len(ids) > 0 {
		if err := u.repository.DeleteImages(ctx, ids); err != nil {
				return nil, err
		}
	}

	return unusedImages, nil
}