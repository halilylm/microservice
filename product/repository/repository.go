package repository

import (
	"context"
	"github.com/halilylm/microservice/domain"
	"time"
)

type ProductRepository interface {
	Insert(ctx context.Context, p *domain.Product) (*domain.Product, error)
	Update(ctx context.Context, p *domain.Product) (*domain.Product, error)
	Delete(ctx context.Context, id int64) error
	GetProductBySlug(ctx context.Context, slug string) (*domain.Product, error)
}

type ProductCacheRepository interface {
	SetProduct(ctx context.Context, key string, expire time.Duration, product *domain.Product) error
	DeleteProduct(ctx context.Context, key string) error
	GetProduct(ctx context.Context, key string) (*domain.Product, error)
}
