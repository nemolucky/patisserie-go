package category_usecase

import "github.com/google/uuid"

type CreateInput struct{
	Title string `json:"title"`
	Parent uuid.UUID `json:"parent"`
}

type DeleteInput struct{
	ID uuid.UUID
}