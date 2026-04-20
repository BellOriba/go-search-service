package database

import (
	"os"

	"github.com/meilisearch/meilisearch-go"
)

func NewMeilisearchClient() meilisearch.ServiceManager {
	host := os.Getenv("MEILI_HOST")
	key := os.Getenv("MEILI_MASTER_KEY")

	return meilisearch.New(
		host,
		meilisearch.WithAPIKey(key),
	)
}

func SetupMeilisearchIndex(client meilisearch.ServiceManager) error {
	index := client.Index("products")

	filterable := []string{"category", "price", "is_featured"}
	sortable := []string{"price", "created_at"}

	iFilterable := make([]interface{}, len(filterable))
	for i, v := range filterable {
		iFilterable[i] = v
	}

	if _, err := index.UpdateFilterableAttributes(&iFilterable); err != nil {
		return err
	}

	if _, err := index.UpdateSortableAttributes(&sortable); err != nil {
		return err
	}

	return nil
}
