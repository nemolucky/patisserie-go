package handler

import (
	"encoding/json"
	"net/http"

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

func (h *UserHandler) Register(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/users", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Get("/me", h.GetMe)
			r.Post("/me/phone/sms", h.SendSMSNewPhone)
			r.Post("/me/phone/verify", h.UpdatePhone)
			r.Delete("/me", h.DeleteAccount)
		})
	})
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {

	user, err := h.usecase.GetMe(r.Context())
	if err != nil {
		http.Error(w, "failed to get user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) SendSMSNewPhone(w http.ResponseWriter, r *http.Request) {

	var input user_usecase.SendSMSNewPhoneInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	err := h.usecase.SendSMSNewPhone(r.Context(), input)
	if err != nil {
		http.Error(w, "failed to send SMS", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) UpdatePhone(w http.ResponseWriter, r *http.Request) {

	var input user_usecase.UpdatePhoneInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	err := h.usecase.UpdatePhone(r.Context(), input)
	if err != nil {
		http.Error(w, "failed to update phone", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	err := h.usecase.DeleteAccount(r.Context())
	if err != nil {
		http.Error(w, "failed to delete account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}