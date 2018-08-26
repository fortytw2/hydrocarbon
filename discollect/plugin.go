package discollect

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// A Plugin is capable of running scrapes, ideally of a common type or against a single site
type Plugin struct {
	Name    string
	Configs []*Config

	// RateLimit is set per-plugin
	RateLimit *RateLimit

	// A ConfigValidator is used to validate dynamically loaded configs and
	// return the title of the config
	ConfigValidator func(*HandlerOpts) (string, error)
	// the Scheduler looks into the past and tells the future
	Scheduler func(*ScheduleRequest) ([]*ScrapeSchedule, error)

	Routes map[string]Handler
}

// Config is a specific configuration of a given plugin
type Config struct {
	// friendly identifier for this config
	Name string
	// Entrypoints is used to start a scrape
	Entrypoints []string
	// DynamicEntry specifies whether this config was created dynamically
	// or is a preset
	DynamicEntry bool
	// Since is used to convey delta information
	Since time.Time
	// Countries is a list of countries this scrape can be executed from
	// in two code, ISO-3166-2 form
	// nil if unused
	Countries []string
}

// Value implements sql.Valuer for config
func (c *Config) Value() (driver.Value, error) {
	j, err := json.Marshal(c)
	return j, err
}

// Scan implements sql.Scanner for config
func (c *Config) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("did not get a []byte from sql driver for *Config")
	}

	return json.Unmarshal(source, c)
}

// HandlerOpts are passed to a Handler
type HandlerOpts struct {
	Config *Config
	// RouteParams are Capture Groups from the Route regexp
	RouteParams []string

	FileStore FileStore

	Client *http.Client
}

// A HandlerResponse is returned from a Handler
type HandlerResponse struct {
	Tasks  []*Task
	Facts  []interface{}
	Errors []error
}

// ErrorResponse is a helper for returning an error from a Handler
func ErrorResponse(err error) *HandlerResponse {
	return &HandlerResponse{
		Errors: []error{
			err,
		},
	}
}

// Response is shorthand for a successful response
func Response(facts []interface{}, tasks ...*Task) *HandlerResponse {
	return &HandlerResponse{
		Facts: facts,
		Tasks: tasks,
	}
}

// A Handler can handle an individual Task
type Handler func(ctx context.Context, ho *HandlerOpts, t *Task) *HandlerResponse

const defaultTimeout = 10 * time.Second

// launchScrape launches a new scrape and enqueues the initial tasks
func launchScrape(ctx context.Context, id uuid.UUID, p *Plugin, cfg *Config, q Queue, ms Metastore) error {
	qts := make([]*QueuedTask, len(cfg.Entrypoints))
	for _, e := range cfg.Entrypoints {
		u, err := uuid.NewRandom()
		if err != nil {
			return err
		}

		qts = append(qts, &QueuedTask{
			Config:   cfg,
			TaskID:   u,
			ScrapeID: id,
			QueuedAt: time.Now(),
			Plugin:   p.Name,
			Retries:  0,
			Task: &Task{
				URL:     e,
				Timeout: defaultTimeout,
			},
		})
	}

	return q.Push(ctx, qts)
}
