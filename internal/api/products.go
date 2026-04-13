package api

import (
	"encoding/json"
	"net/http"

	"github.com/BellOriba/go-search-service/internal/products"
	"github.com/google/uuid"
)

func CreateProductHandler(repo products.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req products.CreateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		product := &products.Product{
			ID: uuid.New(),
			SKU: req.SKU,
			Name: req.Name,
			Slug: req.Slug,
			Description: req.Description,
			Price: req.Price,
			Stock: req.Stock,
			CategoryID: req.CategoryID,
			IsFeatured: req.IsFeatured,
		}

		if err := repo.Create(r.Context(), product); err != nil {
			http.Error(w, "failed to create product", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)
	}
}

