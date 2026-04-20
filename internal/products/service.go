package products

import (
	"context"
	"fmt"
	"strings"
)

type ProductService struct {
	db     Repository
	search SearchRepository
}

func NewProductService(db Repository, search SearchRepository) *ProductService {
	return &ProductService{
		db:     db,
		search: search,
	}
}

func (s *ProductService) Create(ctx context.Context, p *Product) (*Product, error) {
	if err := s.db.Create(ctx, p); err != nil {
		return nil, err
	}

	fullProduct, err := s.db.GetByID(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	if err := s.search.Index(ctx, fullProduct); err != nil {
		return nil, err
	}

	return fullProduct, nil
}

func (s *ProductService) SyncAll(ctx context.Context) (int, error) {
	allProducts, err := s.db.GetAll(ctx)
	if err != nil {
		return 0, err
	}

	if len(allProducts) == 0 {
		return 0, nil
	}

	batchSize := 500
	for i := 0; i < len(allProducts); i += batchSize {
		end := i + batchSize
		if end > len(allProducts) {
			end = len(allProducts)
		}

		chunk := allProducts[i:end]
		if err := s.search.IndexBatch(ctx, chunk); err != nil {
			return i, err
		}
	}

	return len(allProducts), nil
}

func (s *ProductService) SearchProducts(ctx context.Context, query string, category string, maxPrice int64, sortBy, sortOrder string, limit, offset int) ([]ProductIndex, error) {
	var filters []string

	if category != "" {
		filters = append(filters, fmt.Sprintf("category = '%s'", category))
	}
	if maxPrice > 0 {
		filters = append(filters, fmt.Sprintf("price <= %d", maxPrice))
	}

	filterStr := strings.Join(filters, " AND ")

	var sort []string
	if sortBy != "" && sortOrder != "" {
		sort = []string{fmt.Sprintf("%s:%s", sortBy, sortOrder)}
	} else if sortOrder != "" {
		sort = []string{fmt.Sprintf("price:%s", sortOrder)}
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.search.Search(ctx, query, filterStr, sort, limit, offset)
}
