package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/halilylm/microservice/product"
	"github.com/halilylm/microservice/product/repository"
)

const (
	insertQuery    = `INSERT products SET name=?, slug=?, price=?, created_at=now(), updated_at=now()`
	updateQuery    = `UPDATE products SET name=?, price=?, updated_at=now() WHERE id=?`
	deleteQuery    = `DELETE FROM products WHERE id=?`
	getBySlugQuery = `SELECT id, name, slug, price, created_at, updated_at FROM products WHERE slug=?`
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) repository.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Insert(ctx context.Context, p *product.Product) (*product.Product, error) {
	stmt, err := r.db.PrepareContext(ctx, insertQuery)
	if err != nil {
		return nil, err
	}
	res, err := stmt.ExecContext(ctx, p.Name, p.Price)
	if err != nil {
		return nil, err
	}
	p.ID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *productRepository) Update(ctx context.Context, p *product.Product) (*product.Product, error) {
	stmt, err := r.db.PrepareContext(ctx, updateQuery)
	if err != nil {
		return nil, err
	}
	res, err := stmt.ExecContext(ctx, p.Name, p.Price, p.ID)
	if err != nil {
		return nil, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected != 1 {
		return nil, errors.New("no update operated")
	}
	return p, nil
}

func (r *productRepository) Delete(ctx context.Context, id int64) error {
	stmt, err := r.db.PrepareContext(ctx, deleteQuery)
	if err != nil {
		return err
	}
	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.New("no update operated")
	}
	return nil
}

func (r *productRepository) GetProductBySlug(ctx context.Context, slug string) (*product.Product, error) {
	var product product.Product
	if err := r.db.QueryRowContext(ctx, getBySlugQuery, slug).Scan(&product.ID, &product.Name, &product.Slug, &product.Price, &product.CreatedAt, &product.UpdatedAt); err != nil {
		return nil, err
	}
	return &product, nil
}
