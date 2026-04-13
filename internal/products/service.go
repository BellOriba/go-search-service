package products

import "context"

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

func (s *ProductService) Create(ctx context.Context, p *Product) error {
	if err := s.db.Create(ctx, p); err != nil {
		return err
	}

	return s.search.Index(ctx, p)
}

func (s *ProductService) SyncAll(ctx context.Context) (int, error) {
	allProducts, err := s.db.GetAll(ctx)
	if err != nil {
		return 0, err
	}

	for _, p := range allProducts {
		if err := s.search.Index(ctx, p); err != nil {
			return 0, nil
		}
	}

	return len(allProducts), nil
}
