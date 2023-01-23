package server

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type dbPingMock struct {
	err error
}

func (d *dbPingMock) Ping(ctx context.Context) error {
	return d.err
}

func TestHealth(t *testing.T) {
	db := dbPingMock{}
	t.Run("returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		res := httptest.NewRecorder()
		Health(&db)(res, req)
		result := res.Result()
		assert.Equal(t, http.StatusOK, result.StatusCode)
	})
	t.Run("returns 502 when db pings return error", func(t *testing.T) {
		db.err = errors.New("error connecting")
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		res := httptest.NewRecorder()
		Health(&db)(res, req)
		result := res.Result()
		assert.Equal(t, http.StatusBadGateway, result.StatusCode)
	})
}
