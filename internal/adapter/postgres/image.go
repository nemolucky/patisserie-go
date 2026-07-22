package postgres

import (
	"context"
	_ "embed"

	"github.com/google/uuid"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

//go:embed queries/image/create.sql
var CreateImageQuery string

//go:embed queries/image/find_unused.sql
var FindUnusedImagesQuery string

//go:embed queries/image/delete.sql
var DeleteImageQuery string


func (p *Postgres) CreateImage(ctx context.Context, key string) (*uuid.UUID, error) {

	var imageID uuid.UUID

	err := p.Pool.QueryRow(ctx, CreateImageQuery, key).Scan(&imageID)
	if err != nil{
		return nil, err
	}
	
	return &imageID, nil
}

func (p *Postgres) DeleteImages(ctx context.Context, imageIDs []uuid.UUID) (error) {

	_, err := p.Pool.Exec(ctx, DeleteImageQuery, imageIDs)

	return err
}

func (p *Postgres) FindUnusedImages(ctx context.Context) ([]*domain.UnusedImage, error) {

	var unusedImages []*domain.UnusedImage

	rows, err := p.Pool.Query(ctx, FindUnusedImagesQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var image domain.UnusedImage

		if err := rows.Scan(&image.ID, &image.Key); err != nil {
			return nil, err
		}

		unusedImages = append(unusedImages, &image)
	}

	return unusedImages, nil
}