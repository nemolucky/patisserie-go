package worker

import (
	"context"
	"log/slog"
	"time"

	image_usecase "github.com/nemolucky/patisserie-go/internal/usecase/image"
)

type ImageCleanupWorker struct{
	usecase *image_usecase.UseCase
	logger *slog.Logger
	interval time.Duration
}

func NewImageCleanupWorker(
	imageUC *image_usecase.UseCase,
	logger *slog.Logger,
	interval time.Duration,
) *ImageCleanupWorker {
	return &ImageCleanupWorker{
		usecase: imageUC,
		logger: logger,
		interval: interval,
	}
}

func (w *ImageCleanupWorker) Run(ctx context.Context) {
	w.cleanup(ctx)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("Image cleanup worker stopped")
			return

		case <-ticker.C:
			w.cleanup(ctx)
		}
	}
}

func (w *ImageCleanupWorker) cleanup(ctx context.Context) {
	images, err := w.usecase.DeleteUnusedImages(ctx)
	if err != nil {
		w.logger.Error(
			"Cleanup unused images",
			"error", err,
		)
		return
	}

	if len(images) > 0 {
		w.logger.Info(
			"Unused images deleted",
			"count", len(images),
		)
	}
}