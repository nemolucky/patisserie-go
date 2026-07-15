package minio

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/nemolucky/patisserie-go/internal/config"
)

type MinIO struct {
	Client *minio.Client
	Bucket string
	Logger *slog.Logger
}

func New(
	ctx context.Context,
	cfg config.MinioConfig,
	log *slog.Logger,
) (*MinIO, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("check bucket %q: %w", cfg.BucketName, err)
	}

	if !exists {
		if err := client.MakeBucket(
			ctx,
			cfg.BucketName,
			minio.MakeBucketOptions{},
		); err != nil {
			return nil, fmt.Errorf(
				"create bucket %q: %w",
				cfg.BucketName,
				err,
			)
		}

		log.Info(
			"minio bucket created",
			"bucket",
			cfg.BucketName,
		)
	}

	return &MinIO{
		Client: client,
		Bucket: cfg.BucketName,
		Logger: log,
	}, nil
}