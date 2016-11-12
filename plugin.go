package kiasu

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/fortytw2/kiasu/internal/cleanhttp"
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

	// find all configs, with a limit
	Configs func(context.Context, Client, int) ([]Config, error)
	// ensure a configuration is valid
	Validate func(context.Context, Client, Config) error
	// Run launches the given scrape and returns when it is finished
	Run func(context.Context, Client, Config) ([]Post, error)
}

// A Config is a tuple of unique values attached to a scrape
type Config struct {
	InitialURL string
	Since      time.Time
}

// ErrCountryNotFound is returned when a request can't be made in that country
type ErrCountryNotFound struct {
	Country string
}

func (e ErrCountryNotFound) Error() string {
	return fmt.Sprintf("country %s not found in proxy configuration", e.Country)
}

// A Client is used to make all HTTP requests to the outside world
type Client interface {
	Do(*http.Request) (*http.Response, error)
}

// DefaultClient is a simple client, nothing special
func DefaultClient() Client {
	return cleanhttp.DefaultPooledClient()
}
