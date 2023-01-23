package repository

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"github.com/halilylm/microservice/product"
	"sync"
	"time"
)

type MockCacheRepository struct {
	mu       sync.Mutex
	products map[string]*product.Product
}

func NewMockCacheRepository(products map[string]*product.Product) *MockCacheRepository {
	if products == nil {
		products = make(map[string]*product.Product)
	}
	return &MockCacheRepository{products: products}
}

func (mcp *MockCacheRepository) Products() map[string]*product.Product {
	return mcp.products
}

func (mcp *MockCacheRepository) CleanProducts() {
	mcp.products = make(map[string]*product.Product)
}

func (mcp *MockCacheRepository) SetProduct(ctx context.Context, key string, expire time.Duration, product *product.Product) error {
	time.AfterFunc(expire, func() {
		mcp.mu.Lock()
		defer mcp.mu.Unlock()
		delete(mcp.products, key)
	})
	mcp.mu.Lock()
	defer mcp.mu.Unlock()
	mcp.products[key] = product
	return nil
}

func (mcp *MockCacheRepository) DeleteProduct(ctx context.Context, key string) error {
	mcp.mu.Lock()
	defer mcp.mu.Unlock()
	if _, ok := mcp.products[key]; !ok {
		return redis.ErrNil
	}
	delete(mcp.products, key)
	return nil
}

func (mcp *MockCacheRepository) GetProduct(ctx context.Context, key string) (*product.Product, error) {
	mcp.mu.Lock()
	defer mcp.mu.Unlock()
	p, ok := mcp.products[key]
	if !ok {
		return nil, redis.ErrNil
	}
	return p, nil
}
