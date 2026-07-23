package postgres

import (
	"context"
	_ "embed"

	"github.com/google/uuid"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

//go:embed queries/category/create.sql
var CreateCategoryQuery string
//go:embed queries/category/update.sql
var UpdateCategoryQuery string
//go:embed queries/category/delete.sql
var DeleteCategoryQuery string
//go:embed queries/category/list.sql
var ListCategoryQuery string

func (p *Postgres) CreateCategory(ctx context.Context, title string, parent uuid.UUID) (uuid.UUID, error) {

	var categoryID uuid.UUID
	err := p.Pool.QueryRow(ctx, CreateCategoryQuery, title, parent).Scan(&categoryID)

	return categoryID, err
}

func (p *Postgres) UpdateCategory(ctx context.Context, categoryID uuid.UUID, title string) (*domain.Category, error) {

	var category domain.Category

	err := p.Pool.QueryRow(ctx, UpdateCategoryQuery, categoryID, title).Scan(&category.ID, &category.Title, &category.Parent)

	return &category, err
}

func (p *Postgres) DeleteCategory(ctx context.Context, categoryID uuid.UUID) (error) {

	_, err := p.Pool.Exec(ctx, DeleteCategoryQuery, categoryID)

	return err
}

func (p *Postgres) ListCategory(ctx context.Context) ([]domain.Category, error) {

	var categories []domain.Category

	rows, err := p.Pool.Query(ctx, ListCategoryQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var category domain.Category

		if err := rows.Scan(&category.ID, &category.Title); err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}