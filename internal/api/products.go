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

		var productImages []products.ProductImage
		for _, imgReq := range req.Images {
			productImages = append(productImages, products.ProductImage{
				ID:        uuid.New(),
				Path:      imgReq.Path,
				Original:  imgReq.Original,
				Thumbnail: imgReq.Thumbnail,
				IsPrimary: imgReq.IsPrimary,
			})
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
			Images:      productImages,
		}

		fullProduct, err := service.Create(r.Context(), product)
		if err != nil {
			http.Error(w, "failed to create product", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fullProduct)
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

func SearchProductsHandler(service *products.ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		category := r.URL.Query().Get("category")
		sortBy := r.URL.Query().Get("sort_by")
		order := r.URL.Query().Get("order")

		limit := 20
		if l := r.URL.Query().Get("limit"); l != "" {
			fmt.Sscanf(l, "%d", &limit)
		}

		page := 1
		if p := r.URL.Query().Get("page"); p != "" {
			fmt.Sscanf(p, "%d", &page)
		}
		offset := (page - 1) * limit

		var maxPrice int64
		if priceStr := r.URL.Query().Get("max_price"); priceStr != "" {
			fmt.Sscanf(priceStr, "%d", &maxPrice)
		}

		results, err := service.SearchProducts(r.Context(), query, category, maxPrice, sortBy, order, limit, offset)
		if err != nil {
			http.Error(w, "search failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}
