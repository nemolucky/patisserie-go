package handler

import (
	"github.com/go-chi/chi"
	user_usecase "github.com/nemolucky/patisserie-go/internal/usecase/user"
)

type UserHandler struct{
	usecase *user_usecase.UseCase
}

func NewUserHandler(
	usecase *user_usecase.UseCase,
) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func (h *UserHandler) Register(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
	})
}
