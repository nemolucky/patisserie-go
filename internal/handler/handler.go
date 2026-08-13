package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	auth_middleware "github.com/nemolucky/patisserie-go/internal/handler/middleware"
	tokenprovider "github.com/nemolucky/patisserie-go/internal/token_provider"
	auth_usecase "github.com/nemolucky/patisserie-go/internal/usecase/auth"
	category_usecase "github.com/nemolucky/patisserie-go/internal/usecase/category"
	image_usecase "github.com/nemolucky/patisserie-go/internal/usecase/image"
	product_usecase "github.com/nemolucky/patisserie-go/internal/usecase/product"
	user_usecase "github.com/nemolucky/patisserie-go/internal/usecase/user"
)


type Handler struct{
	token *tokenprovider.TokenProvider
	auth *AuthHandler
	product *ProductHandler
	category *CategoryHandler
	image *ImageHandler
	user *UserHandler
}

func New(
	token *tokenprovider.TokenProvider,
	authUC *auth_usecase.UseCase,
	productUC *product_usecase.UseCase,
	categoryUC *category_usecase.UseCase,
	imageUC *image_usecase.UseCase,
	userUC *user_usecase.UseCase,
) *Handler {
	return &Handler{
		token: token,
		auth:  NewAuthHandler(authUC),
		product:  NewProductHandler(productUC),
		category: NewCategoryHandler(categoryUC),
		image:    NewImageHandler(imageUC),
		user:    NewUserHandler(userUC),
	}
}

func (h *Handler) InitRoutes() http.Handler{
	
	r := chi.NewRouter()
	
	r.Use(middleware.Logger, middleware.Recoverer)
	
	authMiddleware := auth_middleware.AuthMiddleware(h.token)

	r.Route("/api/v1", func (r chi.Router){
		h.auth.Register(r)
		h.image.Register(r)
		h.product.Register(r)
		h.category.Register(r)
		h.user.Register(r, authMiddleware)
	})

	return r
}