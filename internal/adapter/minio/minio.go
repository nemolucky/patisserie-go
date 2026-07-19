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

		switch cfg.BucketPolicy {
			case "public":
				if err := setPublicPolicy(ctx, client, cfg.BucketName); err != nil {
					return nil, fmt.Errorf(
						"set bucket %q policy: %w",
						cfg.BucketName,
						err,
					)
				}
			default:
				return nil, fmt.Errorf("unknown bucket policy: %s", cfg.BucketPolicy)
		}

		log.Info(
			"minio bucket created",
			"bucket",
			cfg.BucketName,
			"policy",
			cfg.BucketPolicy,
		)
	}

	return &MinIO{
		Client: client,
		Bucket: cfg.BucketName,
		Logger: log,
	}, nil
}

func setPublicPolicy(ctx context.Context, client *minio.Client, bucket string) error {
	policy := fmt.Sprintf(`{
		"Version":"2012-10-17",
		"Statement":[
			{
				"Effect":"Allow",
				"Principal":{"AWS":["*"]},
				"Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::%s/*"]
			}
		]
	}`, bucket)

	return client.SetBucketPolicy(ctx, bucket, policy)
}

func (m *MinIO) ObjectURL(key string) string {
	return fmt.Sprintf(
		"http://%s/%s/%s",
		m.Client.EndpointURL().Host,
		m.Bucket,
		key,
	)
}