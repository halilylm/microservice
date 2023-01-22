package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	"github.com/halilylm/microservice/domain"
	"github.com/halilylm/microservice/product/repository"
	"time"
)

type productRepository struct {
	conn *redis.Conn
}

func NewProductRepository(conn *redis.Conn) repository.ProductCacheRepository {
	return &productRepository{conn: conn}
}

func (r *productRepository) SetProduct(ctx context.Context, key string, expire time.Duration, product *domain.Product) error {
	productBytes, err := json.Marshal(product)
	if err != nil {
		return err
	}
	if err := r.conn.Set(ctx, key, productBytes, expire).Err(); err != nil {
		return err
	}
	return nil
}

func (r *productRepository) DeleteProduct(ctx context.Context, key string) error {
	if err := r.conn.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

func (r *productRepository) GetProduct(ctx context.Context, key string) (*domain.Product, error) {
	productBytes, err := r.conn.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var product domain.Product
	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, err
	}
	return &product, nil
}
