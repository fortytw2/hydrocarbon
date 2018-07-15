package discollect

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// A Plugin is capable of running scrapes, ideally of a common type or against a single site
type Plugin struct {
	Name     string
	Schedule *Schedule
	Configs  []*Config

	// RateLimit is set per-plugin
	RateLimit *RateLimit

	// A ConfigValidator is used to validate dynamically loaded configs
	ConfigValidator func(*Config) error
	Routes          map[string]Handler
}

// Config is a specific configuration of a given plugin
type Config struct {
	// friendly identifier for this config
	Name string
	// external system ID
	ExternalID string
	// Entrypoints is used to start a scrape
	Entrypoints []string
	// DynamicEntry specifies whether this config was created dynamically
	// or is a preset
	DynamicEntry bool
	// if Entrypoints is null, we can compute the entrypoint (i.e. in a time based Delta)
	ComputedEntry func(ctx context.Context, cfg *Config) error
	// A Plugin can have many types of Scrapes
	// common ones are "full" and "delta"
	Type string
	// Since is used to convey delta information
	Since time.Time
	// Countries is a list of countries this scrape can be executed from
	// in two code, ISO-3166-2 form
	// nil if unused
	Countries []string
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
func launchScrape(ctx context.Context, p *Plugin, cfg *Config, q Queue, ms Metastore) error {
	id, err := ms.StartScrape(ctx, p.Name, cfg)
	if err != nil {
		return err
	}

	if cfg.DynamicEntry {
		if p.ConfigValidator == nil {
			return errors.New("cannot launch DynamicEntry config for plugin without ConfigValidator")
		}
		err = p.ConfigValidator(cfg)
		if err != nil {
			return err
		}
	}

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
