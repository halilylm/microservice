package http

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/halilylm/microservice/domain"
	"github.com/halilylm/microservice/pkg/rest"
	"github.com/halilylm/microservice/product/usecase"
	"net/http"
)

type productHandler struct {
	uc usecase.ProductUseCase
}

func NewProductHandler(uc usecase.ProductUseCase, mux chi.Router) {
	handler := productHandler{uc: uc}
	mux.Post("/api/v1/product", handler.CreateProduct)
}

func (h *productHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(rest.NewBadRequest("validation error"))
		return
	}
	createdProduct, err := h.uc.CreateProduct(r.Context(), &product)
	if err != nil {
		var httpErr *rest.HTTPError
		if errors.As(err, &httpErr) {
			w.WriteHeader(httpErr.Code)
			json.NewEncoder(w).Encode(httpErr)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(rest.NewInternalServerError())
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdProduct)
}
