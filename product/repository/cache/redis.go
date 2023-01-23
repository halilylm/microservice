package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	"github.com/halilylm/microservice/product"
	"github.com/halilylm/microservice/product/repository"
	"time"
)

type productRepository struct {
	client *redis.Client
}

func NewProductRepository(client *redis.Client) repository.ProductCacheRepository {
	return &productRepository{client: client}
}

func (r *productRepository) SetProduct(ctx context.Context, key string, expire time.Duration, product *product.Product) error {
	productBytes, err := json.Marshal(product)
	if err != nil {
		return err
	}
	if err := r.client.Set(ctx, key, productBytes, expire).Err(); err != nil {
		return err
	}
	return nil
}

func (r *productRepository) DeleteProduct(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

func (r *productRepository) GetProduct(ctx context.Context, key string) (*product.Product, error) {
	productBytes, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var product product.Product
	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, err
	}
	return &product, nil
}
