package postgres

import (
	"context"
	_ "embed"

	"github.com/google/uuid"
)

//go:embed queries/image/create.sql
var CreateImageQuery string

func (p *Postgres) CreateImage(ctx context.Context, key string) (*uuid.UUID, error) {

	var imageID uuid.UUID

	err := p.Pool.QueryRow(ctx, CreateImageQuery, key).Scan(&imageID)
	if err != nil{
		return nil, err
	}
	
	return &imageID, nil
}