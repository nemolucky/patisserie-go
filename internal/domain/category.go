package domain

import "github.com/google/uuid"

type Category struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}
