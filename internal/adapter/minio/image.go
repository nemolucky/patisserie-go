package minio

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

func (m *MinIO) Upload(
	ctx context.Context,
	image *domain.UploadImage,
) (*domain.UploadedImage, error) {

	key := fmt.Sprintf(
		"images/%s%s",
		image.ID,
		filepath.Ext(image.Filename),
	)

	reader := bytes.NewReader(image.Data)

	_, err := m.Client.PutObject(
		ctx,
		m.Bucket,
		key,
		reader,
		int64(len(image.Data)),
		minio.PutObjectOptions{
			ContentType: image.ContentType,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("upload image: %w", err)
	}

	return &domain.UploadedImage{
		ID:  image.ID,
		Key: key,
	}, nil
}