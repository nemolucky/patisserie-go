package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/nemolucky/patisserie-go/internal/adapter/minio"
	"github.com/nemolucky/patisserie-go/internal/adapter/postgres"
	"github.com/nemolucky/patisserie-go/internal/config"
	"github.com/nemolucky/patisserie-go/internal/handler"
	image_usecase "github.com/nemolucky/patisserie-go/internal/usecase/image"
	product_usecase "github.com/nemolucky/patisserie-go/internal/usecase/product/crud"
	"github.com/nemolucky/patisserie-go/pkg/logger"
)

func main(){

	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	log := logger.New()
	slog.SetDefault(log)
	
	log.Info("Application started")

	postgres, err := postgres.New(ctx, cfg.Postgres, log)
	if err != nil {
		log.Error("Failed to connect to Postgres", "error", err)
		os.Exit(1)
	}
	defer postgres.Close()
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

	productUC := product_usecase.New(postgres)
	imageUC := image_usecase.New(s3)

	handler := handler.New(productUC, imageUC)
	router := handler.InitRoutes()

	server := &http.Server{
		Addr:    cfg.App.Host + ":" + cfg.App.Port,
		Handler: router,
	}

	log.Info(
		"HTTP server started on address",
		slog.String("addr", server.Addr + "/api/v1"),
	)

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {

		log.Error(
			"HTTP server stopped",
			slog.Any("error", err),
		)
	}
}