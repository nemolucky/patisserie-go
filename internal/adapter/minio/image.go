package minio

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/nemolucky/patisserie-go/internal/domain"
	image_usecase "github.com/nemolucky/patisserie-go/internal/usecase/image"
)

func (m *MinIO) Upload(
	ctx context.Context,
	image image_usecase.UploadImage,
) (*domain.UploadedImage, error) {

	key := fmt.Sprintf(
		"images/%s%s",
		uuid.NewString(),
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
		Key: key,
	}, nil
}

func (m *MinIO) Delete(ctx context.Context, key string) error {
	return m.Client.RemoveObject(ctx, m.Bucket, key, minio.RemoveObjectOptions{})
}