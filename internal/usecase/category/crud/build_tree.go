package category_usecase

import (
	"github.com/google/uuid"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

type CategoryNode struct{
	ID uuid.UUID `json:"id"`
	Title string `json:"title"`
	Children []*CategoryNode `json:"children"`
}

func BuildTree(categories []domain.Category) []*CategoryNode {
	nodes := make(map[uuid.UUID]*CategoryNode, len(categories))

	for _, category := range categories {
		nodes[category.ID] = &CategoryNode{
			ID:    category.ID,
			Title: category.Title,
		}
	}

	var roots []*CategoryNode

	for _, category := range categories {
		node := nodes[category.ID]

		if category.Parent == nil {
			roots = append(roots, node)
			continue
		}

		parent := nodes[*category.Parent]
		parent.Children = append(parent.Children, node)
	}

	return roots
}
