package discollect

import (
	"context"

	"github.com/google/uuid"
)

// A Metastore is used to store the history of all scrape runs
type Metastore interface {
	// StartScrape attempts to start the scrape, returning `true, nil` if the scrape is
	// able to be started
	StartScrape(ctx context.Context, pluginName string, cfg *Config) (id uuid.UUID, err error)
	EndScrape(ctx context.Context, id string, datums, tasks int) error
}

// MemMetastore is a metastore that only stores information in memory
type MemMetastore struct{}

// StartScrape creates an id and starts a scrape in memory
func (MemMetastore) StartScrape(ctx context.Context, pluginName string, cfg *Config) (uuid.UUID, error) {
	return uuid.NewRandom()
}

// EndScrape records the end of a scrape
func (MemMetastore) EndScrape(ctx context.Context, id string, datums, tasks int) error {
	return nil
}
