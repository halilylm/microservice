package http

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/halilylm/microservice/pkg/rest"
	"github.com/halilylm/microservice/product"
	"github.com/halilylm/microservice/product/usecase"
	"net/http"
	"path"
	"strconv"
)

type productHandler struct {
	uc usecase.ProductUseCase
}

func NewProductHandler(uc usecase.ProductUseCase, r chi.Router) {
	handler := productHandler{uc: uc}
	r.Post("/", handler.CreateProduct)
	r.Delete("/", handler.DeleteProduct)
	r.Get("/", handler.GetProductBySlug)
}

func (h *productHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product product.Product
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(rest.NewBadRequest("validation error"))
		return
	}
	validate := validator.New()
	if err := validate.Struct(&product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(rest.NewBadRequest(err.Error()))
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

func (h *productHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)
	pid, err := strconv.Atoi(id)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(rest.NewNotFoundError())
		return
	}
	if err := h.uc.DeleteProduct(r.Context(), int64(pid)); err != nil {
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
	w.WriteHeader(http.StatusOK)
}

func (h *productHandler) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	slug := path.Base(r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	foundProduct, err := h.uc.GetProductBySlug(r.Context(), slug)
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(foundProduct)
}

func (h *productHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product product.Product
	w.Header().Set("Content-Type", "application/json")
	id := path.Base(r.URL.Path)
	pid, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(rest.NewNotFoundError())
		return
	}
	product.ID = int64(pid)
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(rest.NewBadRequest("validation error"))
		return
	}
	validate := validator.New()
	if err := validate.Struct(&product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(rest.NewBadRequest(err.Error()))
		return
	}
	updatedProduct, err := h.uc.UpdateProduct(r.Context(), &product)
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedProduct)
}
