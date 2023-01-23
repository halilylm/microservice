package usecase

import (
	"context"
	"github.com/halilylm/microservice/pkg/rest"
	"github.com/halilylm/microservice/product"
	"github.com/halilylm/microservice/product/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestProductUC_CreateProduct(t *testing.T) {
	t.Parallel()
	repo := repository.NewMockProductRepository(nil)
	cache := repository.NewMockCacheRepository(nil)
	uc := NewProductUC(repo, cache, zap.NewNop())
	t.Run("creates a product", func(t *testing.T) {
		p := product.Product{
			ID:    1,
			Name:  "pear watch",
			Price: 500,
		}
		createdProduct, err := uc.CreateProduct(context.TODO(), &p)
		assert.NoError(t, err)
		assert.NotNil(t, createdProduct)
		assert.Equal(t, 1, len(repo.Products()))
	})
	t.Run("proper slug on collide", func(t *testing.T) {
		repo.CleanProducts()
		p := product.Product{
			ID:    1,
			Name:  "pear watch",
			Price: 500,
		}
		uc.CreateProduct(context.TODO(), &p)
		p2 := product.Product{
			ID:    2,
			Name:  "pear watch",
			Price: 500,
		}
		createdProduct2, err := uc.CreateProduct(context.TODO(), &p2)
		assert.NoError(t, err)
		assert.NotNil(t, createdProduct2)
		assert.Equal(t, "pear-watch-1", createdProduct2.Slug)
		assert.Equal(t, 2, len(repo.Products()))
	})
}

func TestProductUC_DeleteProduct(t *testing.T) {
	t.Parallel()
	repo := repository.NewMockProductRepository(map[int64]*product.Product{
		0: {
			ID:    0,
			Name:  "test",
			Slug:  "test",
			Price: 15,
		},
	})
	cache := repository.NewMockCacheRepository(nil)
	uc := NewProductUC(repo, cache, zap.NewNop())
	t.Run("deletes a product", func(t *testing.T) {
		assert.Equal(t, 1, len(repo.Products()))
		err := uc.DeleteProduct(context.TODO(), 0)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(repo.Products()))
	})
	t.Run("returns error if product not exists", func(t *testing.T) {
		repo.CleanProducts()
		err := uc.DeleteProduct(context.TODO(), 0)
		var httpErr *rest.HTTPError
		assert.ErrorAs(t, err, &httpErr)
	})
}

func TestProductUC_UpdateProduct(t *testing.T) {
	t.Parallel()
	repo := repository.NewMockProductRepository(map[int64]*product.Product{
		0: {
			ID:    0,
			Name:  "test",
			Slug:  "test",
			Price: 15,
		},
	})
	cache := repository.NewMockCacheRepository(nil)
	uc := NewProductUC(repo, cache, zap.NewNop())
	t.Run("updates the product", func(t *testing.T) {
		newProduct := product.Product{
			ID:    0,
			Name:  "lemon",
			Slug:  "test",
			Price: 25,
		}
		updatedProduct, err := uc.UpdateProduct(context.TODO(), &newProduct)
		assert.NoError(t, err)
		assert.Equal(t, "test", newProduct.Slug)
		assert.Equal(t, 1, len(repo.Products()))
		assert.Equal(t, "lemon", updatedProduct.Name)
	})
}

func TestProductUC_GetProductBySlug(t *testing.T) {
	t.Parallel()
	repo := repository.NewMockProductRepository(map[int64]*product.Product{
		0: {
			ID:    0,
			Name:  "test",
			Slug:  "test",
			Price: 15,
		},
	})
	cache := repository.NewMockCacheRepository(nil)
	uc := NewProductUC(repo, cache, zap.NewNop())
	t.Run("returns the product", func(t *testing.T) {
		product, err := uc.GetProductBySlug(context.TODO(), "test")
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, 1, len(repo.Products()))
		assert.Equal(t, 1, len(cache.Products()))
	})
	t.Run("get from the cache", func(t *testing.T) {
		repo.CleanProducts()
		assert.Equal(t, 0, len(repo.Products()))
		product, err := uc.GetProductBySlug(context.TODO(), "test")
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, 1, len(cache.Products()))
	})
	t.Run("returns error if do not exists in both", func(t *testing.T) {
		repo.CleanProducts()
		cache.CleanProducts()
		assert.Equal(t, 0, len(repo.Products()))
		assert.Equal(t, 0, len(cache.Products()))
		product, err := uc.GetProductBySlug(context.TODO(), "test")
		assert.Nil(t, product)
		var httpErr *rest.HTTPError
		assert.ErrorAs(t, err, &httpErr)
	})
}
