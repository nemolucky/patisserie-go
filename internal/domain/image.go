package domain

import "github.com/google/uuid"

type UploadedImage struct {
	Key string `json:"key"`
}

type UnusedImage struct {
	ID  uuid.UUID `json:"id"`
	Key string    `json:"key"`
}