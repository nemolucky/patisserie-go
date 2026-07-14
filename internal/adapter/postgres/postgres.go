package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nemolucky/patisserie-go/internal/config"
)

type Postgres struct {
	Pool   *pgxpool.Pool
	Logger *slog.Logger
}

func New(
	ctx context.Context,
	cfg config.PostgresConfig,
	log *slog.Logger,
) (*Postgres, error) {

	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	log.Info(
		"postgres connected",
		slog.String("host", cfg.Host),
		slog.String("database", cfg.Name),
	)

	return &Postgres{
		Pool:   pool,
		Logger: log,
	}, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}

	p.Logger.Info("postgres connection closed")
}