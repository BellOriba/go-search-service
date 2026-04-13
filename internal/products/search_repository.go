package products

import (
	"context"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
)

type SearchRepository interface {
	Index(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, query string) ([]ProductIndex, error)
}

type meilisearchRepository struct {
	client meilisearch.ServiceManager
	index  string
}

func NewMeilisearchRepository(client meilisearch.ServiceManager) SearchRepository {
	return &meilisearchRepository{
		client: client,
		index:  "products",
	}
}

func (r *meilisearchRepository) Index(ctx context.Context, p *Product) error {
	doc := ProductIndex{
		ID:          p.ID.String(),
		SKU:         p.SKU,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Price:       p.Price,
		Category:    p.CategoryName,
		IsFeatured:  p.IsFeatured,
		CreatedAt:   p.CreatedAt.Unix(),
	}

	_, err := r.client.Index(r.index).AddDocuments([]ProductIndex{doc}, &meilisearch.DocumentOptions{})
	return err
}

func (r *meilisearchRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Index(r.index).DeleteDocument(id, &meilisearch.DocumentOptions{})
	return err
}

func (r *meilisearchRepository) Search(ctx context.Context, query string) ([]ProductIndex, error) {
	searchRes, err := r.client.Index(r.index).Search(query, &meilisearch.SearchRequest{
		Limit: 20,
	})
	if err != nil {
		return nil, err
	}

	var products []ProductIndex
	for _, hit := range searchRes.Hits {
		fmt.Printf("Hit: %v\n", hit)
	}

	return products, nil
}
