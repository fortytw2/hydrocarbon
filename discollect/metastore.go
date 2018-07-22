package discollect

import (
	"context"

	"github.com/google/uuid"
)

type StartedScrape struct {
	ID     uuid.UUID
	FeedID uuid.UUID

	Plugin string
	Config *Config
}

type RunningScrape struct {
	ID     uuid.UUID
	FeedID uuid.UUID

	Plugin string
	Config *Config
}

// A Metastore is used to store the history of all scrape runs and enough meta
// information to allow session resumption on restart of hydrocarbon
type Metastore interface {
	// StartScrapes selects a number of currently STOPPED scrapes, moves them to
	// RUNNING and returns their details
	StartScrapes(ctx context.Context, limit int) ([]*StartedScrape, error)

	// ListScrapes is used to list and filter scrapes, for both session resumption
	// and UI purposes
	ListScrapes(ctx context.Context, statusFilter string) ([]*RunningScrape, error)

	// EndScrape marks a scrape as SUCCESS and records the number of datums and
	// tasks returned
	EndScrape(ctx context.Context, id uuid.UUID, datums, tasks int) error
	// ErrorScrape marks a scrape as ERRORED and adds the error to its list
	ErrorScrape(ctx context.Context, id uuid.UUID, err error) error
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
