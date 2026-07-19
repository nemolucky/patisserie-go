package postgres

import (
	"context"
	_ "embed"
)

//go:embed queries/image/create.sql
var CreateImageQuery string

func (p *Postgres) CreateImage(ctx context.Context, key string) error {

	_, err := p.Pool.Exec(ctx, CreateImageQuery, key)
	
	return err
}