package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	t.Run("returns 200", func(t *testing.T) {
		mux := chi.NewMux()
		Health(mux)
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		res := httptest.NewRecorder()
		mux.ServeHTTP(res, req)
		result := res.Result()
		assert.Equal(t, http.StatusOK, result.StatusCode)
	})
}
