package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
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

	image := &domain.UploadImage{}

	result, err := h.usecase.Upload(
		r.Context(),
		image,
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