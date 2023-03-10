package repository

import (
	"context"
	"database/sql"
	"github.com/halilylm/microservice/product"
	"sync"
)

type MockProductRepository struct {
	sync.Mutex
	products map[int64]*product.Product
}

func NewMockProductRepository(products map[int64]*product.Product) *MockProductRepository {
	if products == nil {
		products = make(map[int64]*product.Product)
	}
	return &MockProductRepository{products: products}
}

func (mpr *MockProductRepository) Products() map[int64]*product.Product {
	return mpr.products
}

func (mpr *MockProductRepository) CleanProducts() {
	mpr.products = make(map[int64]*product.Product)
}

func (mpr *MockProductRepository) Insert(ctx context.Context, p *product.Product) (*product.Product, error) {
	mpr.Lock()
	defer mpr.Unlock()
	if mpr.products == nil {
		mpr.products = make(map[int64]*product.Product)
	}
	if _, ok := mpr.products[p.ID]; ok {
		return nil, sql.ErrNoRows
	}
	mpr.products[p.ID] = p
	return p, nil
}

func (mpr *MockProductRepository) Update(ctx context.Context, p *product.Product) (*product.Product, error) {
	mpr.Lock()
	defer mpr.Unlock()
	if _, ok := mpr.products[p.ID]; !ok {
		return nil, sql.ErrNoRows
	}
	mpr.products[p.ID] = p
	return p, nil
}

func (mpr *MockProductRepository) Delete(ctx context.Context, id int64) error {
	mpr.Lock()
	defer mpr.Unlock()
	if _, ok := mpr.products[id]; !ok {
		return sql.ErrNoRows
	}
	delete(mpr.products, id)
	return nil
}

func (mpr *MockProductRepository) GetProductBySlug(ctx context.Context, slug string) (*product.Product, error) {
	mpr.Lock()
	defer mpr.Unlock()
	for _, v := range mpr.products {
		if v.Slug == slug {
			return v, nil
		}
	}
	return nil, sql.ErrNoRows
}
