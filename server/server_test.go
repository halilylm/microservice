package server_test

import (
	"github.com/halilylm/microservice/test/integration"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestServer_Start(t *testing.T) {
	t.Run("running the server", func(t *testing.T) {
		shutdown := integration.CreateServer()
		defer shutdown()
		req, err := http.Get("http://localhost:9000")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, req.StatusCode)
	})
}
