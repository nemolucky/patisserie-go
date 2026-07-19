package domain

import "github.com/google/uuid"

type UploadImage struct {
	ID          uuid.UUID `json:"id"`
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	Data        []byte `json:"data"`
}

type UploadedImage struct {
	ID  uuid.UUID `json:"id"`
	Key string `json:"key"`
}
