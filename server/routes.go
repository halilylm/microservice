package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	m "github.com/halilylm/microservice/http/middleware"
	"github.com/halilylm/microservice/pkg/database"
	"github.com/halilylm/microservice/product/delivery/http"
	"github.com/halilylm/microservice/product/repository/cache"
	"github.com/halilylm/microservice/product/repository/mysql"
	"github.com/halilylm/microservice/product/usecase"
	"time"
)

func (s *Server) mapRoutes() {
	s.logger.Debug("mapping the routes")
	s.mux.Use(middleware.RequestID)
	s.mux.Use(middleware.RealIP)
	s.mux.Use(m.RequestLogger(s.logger))
	s.mux.Use(middleware.Recoverer)
	db, err := database.NewMysqlConn(database.MysqlConnOptions{
		Host:                  "mysql",
		Port:                  3306,
		User:                  "root",
		Password:              "secret",
		Name:                  "products",
		MaxOpenConnections:    25,
		MaxIdleConnections:    25,
		ConnectionMaxLifetime: 5 * time.Second,
		ConnectionMaxIdleTime: 5 * time.Second,
		Log:                   s.logger,
	})
	if err != nil {
		s.logger.Fatal(err.Error())
	}
	rdb := database.NewRedisConn(database.RedisOptions{
		Host:     "redis",
		Port:     6379,
		Password: "",
		DB:       0,
	})
	s.mux.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/products", func(r chi.Router) {
				crepo := cache.NewProductRepository(rdb.Client)
				prepo := mysql.NewProductRepository(db.DB)
				puc := usecase.NewProductUC(prepo, crepo, s.logger)
				http.NewProductHandler(puc, r)
			})
		})
	})
	s.mux.Get("/health", Health(db, rdb))
}
