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
