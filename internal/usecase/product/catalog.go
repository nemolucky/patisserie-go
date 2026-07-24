package product_usecase

import "context"

func (u *UseCase) ListCatalogProducts(ctx context.Context, input ListCatalogProductsInput) (*ListCatalogProductsOutput, error) {
	
	limit := input.Limit
	offset := input.Offset

	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	products, err := u.repository.ListCatalogProducts(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return &ListCatalogProductsOutput{
		Products: products,
		Limit:    limit,
		Offset:   offset,
	}, nil
}