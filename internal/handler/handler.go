package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	image_usecase "github.com/nemolucky/patisserie-go/internal/usecase/image"
	product_usecase "github.com/nemolucky/patisserie-go/internal/usecase/product/crud"
)


type Handler struct{
	product *ProductHandler
	// categoryHandler *CategoryHandler
	image *ImageHandler
}

func New(
	productUC *product_usecase.UseCase,
	// category *CategoryHandler,
	imageUC *image_usecase.UseCase,
) *Handler {
	return &Handler{
		product:  NewProductHandler(productUC),
		// categoryHandler: category,
		image:    NewImageHandler(imageUC),
	}
}

func (h *Handler) InitRoutes() http.Handler{
	
	r := chi.NewRouter()

	r.Use(middleware.Logger, middleware.Recoverer)

	r.Route("/api/v1", func (r chi.Router){
		h.image.Register(r)
		h.product.Register(r)
	})

	return r
}