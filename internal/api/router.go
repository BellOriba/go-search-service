package api

import (
	"net/http"

	"github.com/BellOriba/go-search-service/internal/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Handler(service *products.ProductService) http.Handler {
	r := chi.NewRouter()

	r.Use(RequestIDMiddleware)
	r.Use(LoggingMiddleware)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/products", CreateProductHandler(service))
		r.Post("/products/sync", SyncProductsHandler(service))
	})

	return r
}
