package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/nemolucky/patisserie-go/internal/adapter/postgres"
	"github.com/nemolucky/patisserie-go/internal/config"
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
		log.Error("failed to connect to postgres", "error", err)
	}
	defer postgres.Close()

	log.Info("Postgres connected")
}