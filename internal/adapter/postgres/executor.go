package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Executor interface{
	Exec(ctx context.Context, sql string, args ...any)(pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any)(pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func ResolveExecutor(ctx context.Context, db Executor) Executor {
	if tx, ok := GetTx(ctx); ok {
		return tx
	}
	return db
}

