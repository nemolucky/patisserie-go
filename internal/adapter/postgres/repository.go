package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nemolucky/patisserie-go/internal/config"
)

type Repository struct {
	pool   *pgxpool.Pool

	Logger *slog.Logger
}

func New(
	ctx context.Context,
	cfg config.PostgresConfig,
	log *slog.Logger,
) (*Repository, error) {

	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("create Database pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping Database: %w", err)
	}

	log.Info(
		"Database connected",
		slog.String("host", cfg.Host),
		slog.String("database", cfg.Name),
	)

	return &Repository{
		pool:   pool,
		Logger: log,
	}, nil
}

func (p *Repository) Close() {
	if p.pool != nil {
		p.pool.Close()
	}

	p.Logger.Info("Database connection closed")
}

func (r *Repository) Executor(ctx context.Context) Executor{
	return ResolveExecutor(ctx, r.pool)
}

func (r *Repository) Exec(ctx context.Context, sql string, args ...any)(pgconn.CommandTag, error){
	return r.Executor(ctx).Exec(ctx, sql, args...)
}

func (r *Repository) Query(ctx context.Context, sql string, args ...any)(pgx.Rows, error){
	return r.Executor(ctx).Query(ctx, sql, args...)
}

func (r *Repository) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row{
	return r.Executor(ctx).QueryRow(ctx, sql, args...)
}