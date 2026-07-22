package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	category_usecase "github.com/nemolucky/patisserie-go/internal/usecase/category/crud"
	image_usecase "github.com/nemolucky/patisserie-go/internal/usecase/image"
	product_usecase "github.com/nemolucky/patisserie-go/internal/usecase/product/crud"
)


type Handler struct{
	product *ProductHandler
	category *CategoryHandler
	image *ImageHandler
}

func New(
	productUC *product_usecase.UseCase,
	categoryUC *category_usecase.UseCase,
	imageUC *image_usecase.UseCase,
) *Handler {
	return &Handler{
		product:  NewProductHandler(productUC),
		category: NewCategoryHandler(categoryUC),
		image:    NewImageHandler(imageUC),
	}
}

func (h *Handler) InitRoutes() http.Handler{
	
	r := chi.NewRouter()

	r.Use(middleware.Logger, middleware.Recoverer)

	r.Route("/api/v1", func (r chi.Router){
		h.image.Register(r)
		h.product.Register(r)
		h.category.Register(r)
	})

	return r
}