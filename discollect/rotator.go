package discollect

import (
	"errors"
	"net"
	"net/http"
	"runtime"
	"time"
)

// Rotator is a proxy rotator interface capable of rotating and rate limiting between many IPs
// TODO(fortytw2): this interface is totally wrong, needs rate limits in it
type Rotator interface {
	Get(c *Config) (*http.Client, error)
}

// DefaultRotator is a no-op rotator that does no proxy rotation
type DefaultRotator struct {
	client *http.Client
}

// NewDefaultRotator provisions a new default rotator
func NewDefaultRotator() *DefaultRotator {
	return &DefaultRotator{
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
			},
		},
	}
}

// Get returns a standard http client
func (dr *DefaultRotator) Get(_ *Config) (*http.Client, error) {
	if dr.client == nil {
		return nil, errors.New("no client configured for DefaultRotator, use WithRotator")
	}
	return dr.client, nil
}
