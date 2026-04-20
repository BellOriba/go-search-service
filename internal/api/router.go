package api

import (
	"net/http"

	"github.com/BellOriba/go-search-service/internal/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Handler(service *products.ProductService, repo *products.PostgresRepository) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
	}))

	r.Use(RequestIDMiddleware)
	r.Use(LoggingMiddleware)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/products/search", SearchProductsHandler(service))
		r.Post("/auth/login", LoginHandler(repo))

		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware)
			r.Post("/products", CreateProductHandler(service))
			r.Post("/products/sync", SyncProductsHandler(service))
		})
	})

	return r
}
