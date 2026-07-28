package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type KeyTx struct{}

func WithTx(ctx context.Context, tx pgx.Tx) context.Context{
	return context.WithValue(ctx, KeyTx{}, tx)
} 

func GetTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(KeyTx{}).(pgx.Tx)
	return tx, ok
}