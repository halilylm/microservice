package http

import (
	"context"
	"encoding/json"
	"github.com/halilylm/microservice/pkg/maps"
	"github.com/halilylm/microservice/pkg/rest"
	"github.com/halilylm/microservice/product"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

type MockProductUsecase struct {
	mu         sync.Mutex
	Products   map[int64]*product.Product
	firstState map[int64]*product.Product
}

func NewMockProductUsecase(products map[int64]*product.Product) *MockProductUsecase {
	if products == nil {
		products = make(map[int64]*product.Product)
	}
	initialState := maps.Copy(products)
	return &MockProductUsecase{
		Products:   products,
		firstState: initialState,
	}
}

func (m *MockProductUsecase) refreshDatabase() {
	m.Products = m.firstState
}

func (m *MockProductUsecase) CreateProduct(ctx context.Context, product *product.Product) (*product.Product, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Products[product.ID] = product
	return product, nil
}

func (m *MockProductUsecase) UpdateProduct(ctx context.Context, product *product.Product) (*product.Product, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.Products[product.ID]; !ok {
		return nil, rest.NewNotFoundError()
	}
	m.Products[product.ID] = product
	return product, nil
}

func (m *MockProductUsecase) DeleteProduct(ctx context.Context, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.Products, id)
	return nil
}

func (m *MockProductUsecase) GetProductBySlug(ctx context.Context, slug string) (*product.Product, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, foundProduct := range m.Products {
		if foundProduct.Slug == slug {
			return foundProduct, nil
		}
	}
	return nil, rest.NewNotFoundError()
}

func TestProductHandler_CreateProduct(t *testing.T) {
	uc := NewMockProductUsecase(nil)
	p := productHandler{uc: uc}
	t.Run("can create a product", func(t *testing.T) {
		product := strings.NewReader(`{"name": "banana watch", "price": 50}`)
		req := httptest.NewRequest(http.MethodPost, "/test", product)
		res := httptest.NewRecorder()
		p.CreateProduct(res, req)
		assert.Equal(t, http.StatusCreated, res.Result().StatusCode)
		assertContentType(t, res, "application/json")
		assert.Equal(t, 1, len(uc.Products))
	})
	t.Run("name is required", func(t *testing.T) {
		uc.refreshDatabase()
		product := strings.NewReader(`{"price": 50}`)
		req := httptest.NewRequest(http.MethodPost, "/test", product)
		res := httptest.NewRecorder()
		p.CreateProduct(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assertContentType(t, res, "application/json")
		assert.Equal(t, 0, len(uc.Products))
	})
	t.Run("price is required", func(t *testing.T) {
		uc.refreshDatabase()
		testProduct := strings.NewReader(`{"name": "banana watch"}`)
		req := httptest.NewRequest(http.MethodPost, "/test", testProduct)
		res := httptest.NewRecorder()
		p.CreateProduct(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assertContentType(t, res, "application/json")
		assert.Equal(t, 0, len(uc.Products))
	})
	t.Run("price should be number", func(t *testing.T) {
		uc.refreshDatabase()
		testProduct := strings.NewReader(`{"name": "banana watch", "price": "abc"}`)
		req := httptest.NewRequest(http.MethodPost, "/test", testProduct)
		res := httptest.NewRecorder()
		p.CreateProduct(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assertContentType(t, res, "application/json")
		assert.Equal(t, 0, len(uc.Products))
	})
	t.Run("returns error for invalid json", func(t *testing.T) {
		uc.refreshDatabase()
		testProduct := strings.NewReader(`{"name": "banana watch", "price": "abc",}`)
		req := httptest.NewRequest(http.MethodPost, "/test", testProduct)
		res := httptest.NewRecorder()
		p.CreateProduct(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assertContentType(t, res, "application/json")
		assert.Equal(t, 0, len(uc.Products))
	})
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	uc := NewMockProductUsecase(map[int64]*product.Product{
		0: {
			ID:    0,
			Name:  "lemon",
			Slug:  "lemon",
			Price: 5,
		},
	})
	p := productHandler{uc: uc}
	t.Run("deletes a product", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/product/0", nil)
		res := httptest.NewRecorder()
		p.DeleteProduct(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assertContentType(t, res, "application/json")
		assert.Equal(t, 0, len(uc.Products))
	})
	t.Run("returns 404 when product not exists", func(t *testing.T) {
		uc.refreshDatabase()
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/product/1", nil)
		res := httptest.NewRecorder()
		p.DeleteProduct(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assertContentType(t, res, "application/json")
		assert.Equal(t, 1, len(uc.Products))
	})
}

func TestProductHandler_GetProductBySlug(t *testing.T) {
	uc := NewMockProductUsecase(map[int64]*product.Product{
		0: {
			ID:    0,
			Name:  "orange book",
			Slug:  "orange-book",
			Price: 500,
		},
	})
	p := productHandler{uc: uc}
	t.Run("gets a product", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/orange-book", nil)
		res := httptest.NewRecorder()
		p.GetProductBySlug(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assertContentType(t, res, "application/json")
	})
	t.Run("returns 404 when product not exists", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/banana-book", nil)
		res := httptest.NewRecorder()
		p.GetProductBySlug(res, req)
		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
		assertContentType(t, res, "application/json")
	})
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	uc := NewMockProductUsecase(map[int64]*product.Product{
		0: {
			ID:    0,
			Name:  "orange book",
			Slug:  "orange-book",
			Price: 500,
		},
	})
	p := productHandler{uc: uc}
	t.Run("updates the product", func(t *testing.T) {
		testProduct := strings.NewReader(`{"name": "banana watch", "price": 200}`)
		req := httptest.NewRequest(http.MethodPut, "/test/0", testProduct)
		res := httptest.NewRecorder()
		p.UpdateProduct(res, req)
		var updatedProduct product.Product
		json.NewDecoder(res.Body).Decode(&updatedProduct)
		assertContentType(t, res, "application/json")
		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, 200, updatedProduct.Price)
		assert.Equal(t, "banana watch", updatedProduct.Name)
		assert.Equal(t, 1, len(uc.Products))
	})
	t.Run("returns 404 when product not exists", func(t *testing.T) {
		testProduct := strings.NewReader(`{"name": "banana watch", "price": 200}`)
		req := httptest.NewRequest(http.MethodPut, "/test/1", testProduct)
		res := httptest.NewRecorder()
		p.UpdateProduct(res, req)
		assertContentType(t, res, "application/json")
		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
		assert.Equal(t, 1, len(uc.Products))
	})
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}
