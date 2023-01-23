package server

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	address string
	mux     chi.Router
	server  *http.Server
	logger  *zap.Logger
}

type Options struct {
	Host   string
	Port   int
	Logger *zap.Logger
}

func New(opts *Options) *Server {
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}
	mux := chi.NewMux()
	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	srv := http.Server{
		Addr:              address,
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
	}
	return &Server{
		address: address,
		mux:     mux,
		server:  &srv,
		logger:  opts.Logger,
	}
}

func (s *Server) Start() error {
	s.mapRoutes()
	s.logger.Info("starting the server at ", zap.String("address", s.address))
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("stopping the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
