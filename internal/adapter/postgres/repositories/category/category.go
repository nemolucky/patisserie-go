package category

import (
	"context"
	_ "embed"

	"github.com/google/uuid"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

//go:embed query/create.sql
var CreateCategoryQuery string
//go:embed query/update.sql
var UpdateCategoryQuery string
//go:embed query/delete.sql
var DeleteCategoryQuery string
//go:embed query/list.sql
var ListCategoryQuery string

func (r *CategoryRepository) CreateCategory(ctx context.Context, title string, parent uuid.UUID) (uuid.UUID, error) {

	var categoryID uuid.UUID
	err := r.QueryRow(ctx, CreateCategoryQuery, title, parent).Scan(&categoryID)

	return categoryID, err
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, categoryID uuid.UUID, title string) (*domain.Category, error) {

	var category domain.Category

	err := r.QueryRow(ctx, UpdateCategoryQuery, categoryID, title).Scan(&category.ID, &category.Title, &category.Parent)

	return &category, err
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, categoryID uuid.UUID) (error) {

	_, err := r.Exec(ctx, DeleteCategoryQuery, categoryID)

	return err
}

func (r *CategoryRepository) ListCategory(ctx context.Context) ([]domain.Category, error) {

	var categories []domain.Category

	rows, err := r.Query(ctx, ListCategoryQuery)
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