package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BellOriba/go-search-service/internal/products"
	"github.com/google/uuid"
)

func CreateProductHandler(service *products.ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req products.CreateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errorMessage := fmt.Sprintf("failed to create product: %v", err)
			http.Error(w, errorMessage, http.StatusBadRequest)
			return
		}

		product := &products.Product{
			ID:          uuid.New(),
			SKU:         req.SKU,
			Name:        req.Name,
			Slug:        req.Slug,
			Description: req.Description,
			Price:       req.Price,
			Stock:       req.Stock,
			CategoryID:  req.CategoryID,
			IsFeatured:  req.IsFeatured,
		}

		if err := service.Create(r.Context(), product); err != nil {
			http.Error(w, "failed to create product", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)
	}
}

func SyncProductsHandler(service *products.ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count, err := service.SyncAll(r.Context())
		if err != nil {
			http.Error(w, "failed to sync products", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "sync completed",
			"count":   count,
		})
	}
}
