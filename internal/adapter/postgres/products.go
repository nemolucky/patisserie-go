package postgres

import (
	"context"

	_ "embed"

	"github.com/nemolucky/patisserie-go/internal/domain"
)

//go:embed queries/products/list_catalog.sql
var ListCatalogProductsQuery string


func (p *Postgres) ListCatalogProducts(ctx context.Context, limit int, offset int) ([]*domain.CatalogProduct, error) {

	products := make([]*domain.CatalogProduct, 0, limit)

	rows, err := p.Pool.Query(ctx, ListCatalogProductsQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var product domain.CatalogProduct

		if err := rows.Scan(&product.ID, &product.Title, &product.Price, &product.Image.URL, &product.Image.SortOrder, &product.Categories); err != nil {
			return nil, err
		}

		products = append(products, &product)
	}
	if err := rows.Err(); err != nil {
    return nil, err
	}


	return products, nil
}