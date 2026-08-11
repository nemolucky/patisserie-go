package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	auth_usecase "github.com/nemolucky/patisserie-go/internal/usecase/auth"
)

var maxAge int = 60 * 60 * 24 * 7

type AuthHandler struct{
	usecase *auth_usecase.UseCase
}

func NewAuthHandler(
	usecase *auth_usecase.UseCase,
) *AuthHandler {
	return &AuthHandler{
		usecase: usecase,
	}
}

func (h *AuthHandler) Register(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/sms", h.SendSMSCode)
		r.Post("/verify", h.VerifySMSCode)
		r.Post("/register", h.RegisterUser)
		r.Post("/refresh", h.Refresh)
		r.Post("/logout", h.Logout)
	})
}

func (h *AuthHandler) SendSMSCode(w http.ResponseWriter, r *http.Request) {
	
	body := r.Body
	defer r.Body.Close()

	var input auth_usecase.SendSMSCodeInput
	if err := json.NewDecoder(body).Decode(&input); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	if err := h.usecase.SendSMSCode(r.Context(), input); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) VerifySMSCode(w http.ResponseWriter, r *http.Request) {
	
	body := r.Body
	defer r.Body.Close()

	var input auth_usecase.VerifySMSCodeInput
	if err := json.NewDecoder(body).Decode(&input); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	output, err := h.usecase.VerifySMSCode(r.Context(), input)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	if output.RegistrationRequired == true{
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(output.RegistrationRequired); err != nil {
			http.Error(
				w,
				"failed to encode response",
				http.StatusInternalServerError,
			)
			return
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
    Name:     "refresh_token",
    Value:    output.TokenPair.RefreshToken,
    HttpOnly: true,
    Secure:   true,
    SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge: maxAge,
	})

	response := struct{
		RegistrationRequired bool `json:"registration_required"`
		AccessToken string `json:"access_token"`
	}{
		RegistrationRequired: output.RegistrationRequired,
		AccessToken: output.TokenPair.AccessToken,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
		return
	}
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	
	body := r.Body
	defer r.Body.Close()

	var input auth_usecase.RegisterUserInput
	if err := json.NewDecoder(body).Decode(&input); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	output, err := h.usecase.RegisterUser(r.Context(), input)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	http.SetCookie(w, &http.Cookie{
    Name:     "refresh_token",
    Value:    output.RefreshToken,
    HttpOnly: true,
    Secure:   true,
    SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge: maxAge,
	})

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(struct{
		AccessToken string `json:"access_token"`
	}{
		AccessToken: output.AccessToken,
	},
	); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
		return
	}
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	output, err := h.usecase.Refresh(r.Context(), cookie.Value)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    output.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge: maxAge,
	})

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(struct{
		AccessToken string `json:"access_token"`
	}{
		AccessToken: output.AccessToken,
	},
	); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
		return
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusNoContent)
}