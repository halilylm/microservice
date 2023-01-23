package repository

import (
	"context"
	"github.com/halilylm/microservice/product"
	"time"
)

type ProductRepository interface {
	Insert(ctx context.Context, p *product.Product) (*product.Product, error)
	Update(ctx context.Context, p *product.Product) (*product.Product, error)
	Delete(ctx context.Context, id int64) error
	GetProductBySlug(ctx context.Context, slug string) (*product.Product, error)
}

type ProductCacheRepository interface {
	SetProduct(ctx context.Context, key string, expire time.Duration, product *product.Product) error
	DeleteProduct(ctx context.Context, key string) error
	GetProduct(ctx context.Context, key string) (*product.Product, error)
}
