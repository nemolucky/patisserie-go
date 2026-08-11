package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nemolucky/patisserie-go/internal/adapter/minio"
	"github.com/nemolucky/patisserie-go/internal/adapter/postgres"
	"github.com/nemolucky/patisserie-go/internal/adapter/postgres/repositories/category"
	"github.com/nemolucky/patisserie-go/internal/adapter/postgres/repositories/image"
	"github.com/nemolucky/patisserie-go/internal/adapter/postgres/repositories/products"
	"github.com/nemolucky/patisserie-go/internal/adapter/postgres/repositories/users"
	"github.com/nemolucky/patisserie-go/internal/adapter/redis"
	smsprovider "github.com/nemolucky/patisserie-go/internal/adapter/sms_provider"
	"github.com/nemolucky/patisserie-go/internal/config"
	"github.com/nemolucky/patisserie-go/internal/handler"
	tokenprovider "github.com/nemolucky/patisserie-go/internal/token_provider"

	auth_usecase "github.com/nemolucky/patisserie-go/internal/usecase/auth"
	category_usecase "github.com/nemolucky/patisserie-go/internal/usecase/category"
	image_usecase "github.com/nemolucky/patisserie-go/internal/usecase/image"
	product_usecase "github.com/nemolucky/patisserie-go/internal/usecase/product"

	"github.com/nemolucky/patisserie-go/internal/worker"
	"github.com/nemolucky/patisserie-go/pkg/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	log := logger.New()
	slog.SetDefault(log)

	log.Info("Application started")

	repo, err := postgres.New(ctx, cfg.Postgres, log)
	if err != nil {
		log.Error(
			"Failed to connect Postgres",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	defer repo.Close()

	userRepo := users.New(repo)
	productRepo := products.New(repo)
	categoryRepo := category.New(repo)
	imageRepo := image.New(repo)

	tm := postgres.NewTransactionManager(repo)

	log.Info("Postgres connected")

	s3, err := minio.New(
		ctx,
		cfg.Minio,
		log,
	)
	if err != nil {
		log.Error(
			"Failed initialize MinIO",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	log.Info(
		"MinIO connected",
		slog.String("endpoint", cfg.Minio.Endpoint),
		slog.String("bucket", cfg.Minio.BucketName),
	)

	rds := redis.New(ctx, cfg.Redis, log)

	smsProvider := &smsprovider.SMSProvider{}

	tokenProvider := tokenprovider.NewTokenProvider(cfg.JWT)

	authUC := auth_usecase.New(userRepo, rds, smsProvider, tokenProvider)
	productUC := product_usecase.New(productRepo, tm)
	categoryUC := category_usecase.New(categoryRepo)
	imageUC := image_usecase.New(s3, imageRepo)

	workers := []worker.Worker{
		worker.NewImageCleanupWorker(
			imageUC,
			log,
			time.Hour,
		),
	}

	var wg sync.WaitGroup

	for _, w := range workers {
		wg.Add(1)

		go func(w worker.Worker) {
			defer wg.Done()
			w.Run(ctx)
		}(w)
	}

	handler := handler.New(
		tokenProvider,
		authUC,
		productUC,
		categoryUC,
		imageUC,
	)

	router := handler.InitRoutes()

	server := &http.Server{
		Addr:    cfg.App.Host + ":" + cfg.App.Port,
		Handler: router,
	}

	go func() {
		<-ctx.Done()

		log.Info("Shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error(
				"Failed shutdown HTTP server",
				slog.Any("error", err),
			)
		}
	}()

	log.Info(
		"HTTP server started",
		slog.String(
			"addr",
			server.Addr+"/api/v1",
		),
	)

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {

		log.Error(
			"HTTP server stopped",
			slog.Any("error", err),
		)
	}

	stop()

	wg.Wait()

	log.Info("Application stopped")
}