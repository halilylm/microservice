package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/halilylm/microservice/pkg/rest"
	"github.com/halilylm/microservice/product"
	"github.com/halilylm/microservice/product/repository"
	"go.uber.org/zap"
	"time"
)

type productUC struct {
	repo   repository.ProductRepository
	cache  repository.ProductCacheRepository
	logger *zap.Logger
}

func NewProductUC(repo repository.ProductRepository, cache repository.ProductCacheRepository, logger *zap.Logger) ProductUseCase {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &productUC{repo: repo, cache: cache, logger: logger}
}

func (p *productUC) CreateProduct(ctx context.Context, product *product.Product) (*product.Product, error) {
	genSlug := slug.Make(product.Name)
	for i := 1; ; i++ {
		product, _ := p.repo.GetProductBySlug(ctx, genSlug)
		if product == nil {
			break
		}
		genSlug = fmt.Sprintf("%s-%d", genSlug, i)
	}
	product.Slug = genSlug
	createdProduct, err := p.repo.Insert(ctx, product)
	if err != nil {
		return nil, rest.NewInternalServerError()
	}
	return createdProduct, nil
}

func (p *productUC) UpdateProduct(ctx context.Context, product *product.Product) (*product.Product, error) {
	updatedProduct, err := p.repo.Update(ctx, product)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, rest.NewNotFoundError()
		}
		return nil, rest.NewInternalServerError()
	}
	return updatedProduct, nil
}

func (p *productUC) DeleteProduct(ctx context.Context, id int64) error {
	if err := p.repo.Delete(ctx, id); err != nil {
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return rest.NewNotFoundError()
			}
			return rest.NewInternalServerError()
		}
	}
	return nil
}

func (p *productUC) GetProductBySlug(ctx context.Context, slug string) (*product.Product, error) {
	// check if exists on the cache
	if foundProduct, err := p.cache.GetProduct(ctx, slug); err == nil {
		p.logger.Debug("getting product from the cache", zap.Int64("id", foundProduct.ID))
		return foundProduct, nil
	}
	// get product in sql
	foundProduct, err := p.repo.GetProductBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, rest.NewNotFoundError()
		}
		return nil, rest.NewInternalServerError()
	}
	// store it in the cache
	if err := p.cache.SetProduct(ctx, slug, 10*time.Second, foundProduct); err != nil {
		p.logger.Error("could not cache the product", zap.Int64("id", foundProduct.ID), zap.Error(err))
	}
	return foundProduct, nil
}

type ProductUseCase interface {
	CreateProduct(ctx context.Context, product *product.Product) (*product.Product, error)
	UpdateProduct(ctx context.Context, product *product.Product) (*product.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
	GetProductBySlug(ctx context.Context, slug string) (*product.Product, error)
}
