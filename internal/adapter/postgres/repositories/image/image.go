package image

import (
	"context"
	_ "embed"

	"github.com/google/uuid"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

//go:embed query/create.sql
var CreateImageQuery string

//go:embed query/find_unused.sql
var FindUnusedImagesQuery string

//go:embed query/delete.sql
var DeleteImageQuery string


func (r *ImageRepository) CreateImage(ctx context.Context, key string) (*uuid.UUID, error) {

	var imageID uuid.UUID

	err := r.QueryRow(ctx, CreateImageQuery, key).Scan(&imageID)
	if err != nil{
		return nil, err
	}
	
	return &imageID, nil
}

func (r *ImageRepository) DeleteImages(ctx context.Context, imageIDs []uuid.UUID) (error) {

	_, err := r.Exec(ctx, DeleteImageQuery, imageIDs)

	return err
}

func (r *ImageRepository) FindUnusedImages(ctx context.Context) ([]*domain.UnusedImage, error) {

	var unusedImages []*domain.UnusedImage

	rows, err := r.Query(ctx, FindUnusedImagesQuery)
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