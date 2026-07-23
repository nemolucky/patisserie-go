package category_usecase

import "github.com/google/uuid"

type CreateInput struct{
	Title string `json:"title"`
	Parent uuid.UUID `json:"parent"`
}

type UpdateInput struct{
	Title string `json:"title"`
}

type DeleteInput struct{
	ID uuid.UUID `json:"id"`
}