package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	category_usecase "github.com/nemolucky/patisserie-go/internal/usecase/category"
)

type CategoryHandler struct{
	usecase *category_usecase.UseCase
}

func NewCategoryHandler(
	usecase *category_usecase.UseCase,
) *CategoryHandler {
	return &CategoryHandler{
		usecase: usecase,
	}
}

func (h *CategoryHandler) Register(r chi.Router) {
	r.Route("/categories", func(r chi.Router) {
		r.Get("/", h.List)
		r.Patch("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
		r.Get("/tree", h.Tree)
		r.Post("/", h.Create)
	})
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	
	categories, err := h.usecase.List(r.Context())
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
}

func (h *CategoryHandler) Tree(w http.ResponseWriter, r *http.Request) {
	
	categories, err := h.usecase.Tree(r.Context())
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	
	body := r.Body
	defer r.Body.Close()

	var input category_usecase.CreateInput
	if err := json.NewDecoder(body).Decode(&input); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
	}

	categoryID, err := h.usecase.Create(r.Context(), input)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(categoryID.String()))
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(
			w,
			"invalid id",
			http.StatusBadRequest,
		)
		return
	}

	body := r.Body
	defer r.Body.Close()

	var input category_usecase.UpdateInput
	if err := json.NewDecoder(body).Decode(&input); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
	}

	output, err := h.usecase.Update(r.Context(), id, input)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(
			w,
			"invalid id",
			http.StatusBadRequest,
		)
		return
	}

	input := category_usecase.DeleteInput{
		ID: id,
	}

	if err := h.usecase.Delete(r.Context(), input); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}