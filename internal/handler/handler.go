package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	image_usecase "github.com/nemolucky/patisserie-go/internal/usecase/image"
)


type Handler struct{
	// productHandler *ProductHandler
	// categoryHandler *CategoryHandler
	image *ImageHandler
}

func New(
	// product *ProductHandler,
	// category *CategoryHandler,
	imageUC *image_usecase.UseCase,
) *Handler {
	return &Handler{
		// productHandler:  product,
		// categoryHandler: category,
		image:    NewImageHandler(imageUC),
	}
}

func (h *Handler) InitRoutes() http.Handler{
	
	r := chi.NewRouter()

	r.Use(middleware.Logger, middleware.Recoverer)

	h.image.Register(r)

	return r
}