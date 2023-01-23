package cache

import (
	"context"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/halilylm/microservice/product"
	"github.com/halilylm/microservice/product/repository"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func setupRedis(t *testing.T) repository.ProductCacheRepository {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	return NewProductRepository(client)
}

func TestProductRepository_SetProduct(t *testing.T) {
	productRepo := setupRedis(t)
	key := uuid.NewString()
	newProduct := product.Product{
		Name:  "lemon",
		Slug:  "lemon",
		Price: 15,
	}
	err := productRepo.SetProduct(context.TODO(), key, 10*time.Second, &newProduct)
	assert.NoError(t, err)
	assert.Nil(t, err)
}

func TestProductRepository_DeleteProduct(t *testing.T) {
	productRepo := setupRedis(t)
	key := uuid.NewString()
	err := productRepo.DeleteProduct(context.TODO(), key)
	assert.NoError(t, err)
	assert.Nil(t, err)
}

func TestProductRepository_GetProduct(t *testing.T) {
	productRepo := setupRedis(t)
	key := uuid.NewString()
	product := product.Product{
		Name:  "banana watch",
		Slug:  "banana-watch",
		Price: 50,
	}
	err := productRepo.SetProduct(context.TODO(), key, 10*time.Second, &product)
	assert.NoError(t, err)
	assert.Nil(t, err)
	prod, err := productRepo.GetProduct(context.TODO(), key)
	assert.NoError(t, err)
	assert.NotNil(t, prod)
}
