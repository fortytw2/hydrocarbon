package discollect

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Scrape struct {
	ID     uuid.UUID `json:"id"`
	FeedID uuid.UUID `json:"feed_id"`

	CreatedAt        time.Time `json:"created_at"`
	ScheduledStartAt time.Time `json:"scheduled_start_at"`
	StartedAt        time.Time `json:"started_at"`
	EndedAt          time.Time `json:"ended_at"`

	State  string   `json:"state"`
	Errors []string `json:"errors"`

	TotalDatums  int `json:"total_datums"`
	TotalRetries int `json:"total_retries"`
	TotalTasks   int `json:"total_tasks"`

	Plugin string  `json:"plugin"`
	Config *Config `json:"config"`
}

// A Metastore is used to store the history of all scrape runs and enough meta
// information to allow session resumption on restart of hydrocarbon
type Metastore interface {
	// StartScrapes selects a number of currently STOPPED scrapes, moves them to
	// RUNNING and returns their details
	StartScrapes(ctx context.Context, limit int) ([]*Scrape, error)

	// ListScrapes is used to list and filter scrapes, for both session resumption
	// and UI purposes
	ListScrapes(ctx context.Context, statusFilter string, limit, offset int) ([]*Scrape, error)

	// ScheduleForwardScrapes adds scrapes that should be run to the future set
	ScheduleForwardScrapes(ctx context.Context, limit int) error

	// EndScrape marks a scrape as SUCCESS and records the number of datums and
	// tasks returned
	EndScrape(ctx context.Context, id uuid.UUID, datums, retries, tasks int) error
	// ErrorScrape marks a scrape as ERRORED and adds the error to its list
	ErrorScrape(ctx context.Context, id uuid.UUID, err error) error
}

// MemMetastore is a metastore that only stores information in memory
// TODO: allow this to function again.
type MemMetastore struct{}
