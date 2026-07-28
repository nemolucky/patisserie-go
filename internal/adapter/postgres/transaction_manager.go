package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)


type TransactionManager struct{
	*Repository
}

func NewTransactionManager(db *Repository) *TransactionManager{
	return &TransactionManager{
		db,
	}
}

func (tm *TransactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil{
		return err
	}
	defer tx.Rollback(ctx)

	ctx = WithTx(ctx, tx)

	err = fn(ctx)
	if err != nil{
		return err
	}

	return tx.Commit(ctx)
}