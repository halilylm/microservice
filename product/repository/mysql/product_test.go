package mysql

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/halilylm/microservice/product"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestProductRepository_Insert(t *testing.T) {
	newProduct := product.Product{
		Name:  "watch",
		Price: 15,
	}
	db, mock := createMockDB(t)
	defer func() {
		_ = db.Close()
	}()
	prep := mock.ExpectPrepare(insertQuery)
	prep.ExpectExec().WithArgs(newProduct.Name, newProduct.Price).WillReturnResult(sqlmock.NewResult(1, 1))
	p := NewProductRepository(db)
	prod, err := p.Insert(context.TODO(), &newProduct)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, prod.ID)
}

func TestProductRepository_GetProductBySlug(t *testing.T) {
	db, mock := createMockDB(t)
	defer func() {
		_ = db.Close()
	}()
	createdAt := time.Now()
	updatedAt := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "slug", "price", "created_at", "updated_at"}).
		AddRow(1, "red lemon", "red-lemon", 5, createdAt, updatedAt)
	mock.ExpectQuery(getBySlugQuery).WillReturnRows(rows)
	p := NewProductRepository(db)
	slug := "red-lemon"
	prod, err := p.GetProductBySlug(context.TODO(), slug)
	assert.NoError(t, err)
	assert.NotNil(t, prod)
}

func TestProductRepository_Delete(t *testing.T) {
	db, mock := createMockDB(t)
	defer func() {
		_ = db.Close()
	}()
	prep := mock.ExpectPrepare(deleteQuery)
	prep.ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	p := NewProductRepository(db)
	err := p.Delete(context.TODO(), 1)
	assert.NoError(t, err)
}

func TestProductRepository_Update(t *testing.T) {
	product := product.Product{
		ID:    1,
		Name:  "banana",
		Price: 5,
	}
	db, mock := createMockDB(t)
	defer func() {
		_ = db.Close()
	}()
	prep := mock.ExpectPrepare(updateQuery)
	prep.ExpectExec().WithArgs(product.Name, product.Price, product.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	p := NewProductRepository(db)
	updateProduct, err := p.Update(context.TODO(), &product)
	assert.NoError(t, err)
	assert.NotNil(t, updateProduct)
}

func createMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an rest '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}
