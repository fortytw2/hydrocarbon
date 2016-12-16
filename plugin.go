package hydrocarbon

import (
	"errors"
	"net/http"
	"time"

	"github.com/fortytw2/hydrocarbon/internal/cleanhttp"
)

// Plugin Errors
var (
	ErrPluginUninitialized = errors.New("plugin not initialized")
)

// An Instantiator is used to create new, stateful copies of a plugin
type Instantiator func() (*Plugin, error)

// A Plugin scrapes something and returns something
type Plugin struct {
	Name string

	// find all configs, and paginate through them
	// return configs found, max configs, error
	Configs func(Client, *Pagination) ([]Config, int, error)
	// ensure a configuration is valid
	Validate func(Client, Config) error
	// Run launches the given scrape and returns when it is finished
	Run func(Client, Config) ([]Post, error)
}

// A Config is a tuple of unique values attached to a scrape
type Config struct {
	InitialURL string
	Since      time.Time
}

// A Client is used to make all HTTP requests to the outside world
type Client interface {
	Do(*http.Request) (*http.Response, error)
}

// DefaultClient is a simple client, nothing special
func DefaultClient() Client {
	return cleanhttp.DefaultPooledClient()
}
