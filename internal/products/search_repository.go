package products

import (
	"context"
	"encoding/json"

	"github.com/meilisearch/meilisearch-go"
)

type SearchRepository interface {
	Index(ctx context.Context, p *Product) error
	IndexBatch(ctx context.Context, products []*Product) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, query string, filter string, sort []string, limit, offset int) ([]ProductIndex, error)
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
	imageURL := ""
	if len(p.Images) > 0 {
		imageURL = p.Images[0].Original
	}

	doc := ProductIndex{
		ID:          p.ID.String(),
		SKU:         p.SKU,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Price:       p.Price,
		Category:    p.CategoryName,
		Image:       imageURL,
		IsFeatured:  p.IsFeatured,
		CreatedAt:   p.CreatedAt.Unix(),
	}

	_, err := r.client.Index(r.index).AddDocuments([]ProductIndex{doc}, &meilisearch.DocumentOptions{})
	return err
}

func (r *meilisearchRepository) IndexBatch(ctx context.Context, products []*Product) error {
	if len(products) == 0 {
		return nil
	}

	docs := make([]ProductIndex, len(products))
	for i, p := range products {
		imageURL := ""
		if len(p.Images) > 0 { imageURL = p.Images[0].Original }

		docs[i] = ProductIndex{
			ID:          p.ID.String(),
			SKU:         p.SKU,
			Name:        p.Name,
			Slug:        p.Slug,
			Description: p.Description,
			Price:       p.Price,
			Category:    p.CategoryName,
			Image: imageURL,
			IsFeatured:  p.IsFeatured,
			CreatedAt:   p.CreatedAt.Unix(),
		}
	}

	_, err := r.client.Index(r.index).AddDocuments(docs, &meilisearch.DocumentOptions{})
	return err
}

func (r *meilisearchRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Index(r.index).DeleteDocument(id, &meilisearch.DocumentOptions{})
	return err
}

func (r *meilisearchRepository) Search(ctx context.Context, query string, filter string, sort []string, limit, offset int) ([]ProductIndex, error) {
	searchRes, err := r.client.Index(r.index).Search(query, &meilisearch.SearchRequest{
		Limit:  int64(limit),
		Offset: int64(offset),
		Filter: filter,
		Sort:   sort,
	})
	if err != nil {
		return nil, err
	}

	var products []ProductIndex
	for _, hit := range searchRes.Hits {
		var p ProductIndex
		jsonData, _ := json.Marshal(hit)
		json.Unmarshal(jsonData, &p)
		products = append(products, p)
	}

	return products, nil
}
