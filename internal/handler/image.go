package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/nemolucky/patisserie-go/internal/domain"
	image_usecase "github.com/nemolucky/patisserie-go/internal/usecase/image"
)

type ImageHandler struct{
	usecase *image_usecase.UseCase
}

func NewImageHandler(
	usecase *image_usecase.UseCase,
) *ImageHandler {
	return &ImageHandler{
		usecase: usecase,
	}
}

func (h *ImageHandler) Register(r chi.Router) {
	r.Route("/images", func(r chi.Router) {
		r.Post("/", h.Upload)
	})
}

func (h *ImageHandler) Upload(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(50 << 20)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	headers := r.MultipartForm.File["images"]

	if len(headers) == 0 {
    http.Error(
        w,
        "images are required",
        http.StatusBadRequest,
    )
    return
	}

	images := make([]*domain.UploadImage, 0, len(headers))

	for _, header := range headers {
		
		if header.Size > 10<<20 {
			http.Error(
				w,
				"image too big",
				http.StatusBadRequest,
			)
			return	
		}

		file, err := header.Open()
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		images = append(images, &domain.UploadImage{
				ID:          uuid.New(),
				Filename:    header.Filename,
				ContentType: header.Header.Get("Content-Type"),
				Data:        data,
		})
	}

	result, err := h.usecase.Upload(
    r.Context(),
    images,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(result)
}