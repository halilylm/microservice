package server

import (
	"context"
	"encoding/json"
	"github.com/halilylm/microservice/pkg/rest"
	"net/http"
)

type pinger interface {
	Ping(ctx context.Context) error
}

func Health(pingers ...pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, p := range pingers {
			if err := p.Ping(context.TODO()); err != nil {
				restErr := rest.NewStatusBadGateway()
				w.WriteHeader(restErr.Code)
				json.NewEncoder(w).Encode(restErr)
				return
			}
		}
	}
}
