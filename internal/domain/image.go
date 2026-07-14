package domain

type UploadImage struct {
	ID          string `json:"id"`
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	Data        []byte `json:"data"`
}

type UploadedImage struct {
	ID  string `json:"id"`
	Key string `json:"key"`
	URL string `json:"url"`
}
