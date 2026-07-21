package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	product_usecase "github.com/nemolucky/patisserie-go/internal/usecase/product/crud"
)

type ProductHandler struct{
	usecase *product_usecase.UseCase
}

func NewProductHandler(
	usecase *product_usecase.UseCase,
) *ProductHandler {
	return &ProductHandler{
		usecase: usecase,
	}
}

func (h *ProductHandler) Register(r chi.Router) {
	r.Route("/products", func(r chi.Router) {
		r.Get("/", h.ListCatalogProducts)
		r.Post("/", h.CreateProduct) 
	})
}

func (h *ProductHandler) ListCatalogProducts(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	limitStr := query.Get("limit")
	offsetStr := query.Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	input := product_usecase.ListCatalogProductsInput{
		Limit:  limit,
		Offset: offset,
	}

	output, err := h.usecase.ListCatalogProducts(
		r.Context(),
		input,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {

	body := r.Body
	defer r.Body.Close()

	var input product_usecase.CreateProductInput
	if err := json.NewDecoder(body).Decode(&input); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
	}

	productID, err := h.usecase.Create(r.Context(), input)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(productID.String()))
}