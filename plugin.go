package kiasu

import "time"

// A Plugin is used to update and refresh feeds
type Plugin interface {
	CanCheckFeed(*Feed) error
	CheckFeed(f *Feed, since time.Time) ([]Post, error)
	Healthcheck() (*Healthcheck, error)
}

// A Healthcheck is a health check of a plugin (is it running)
type Healthcheck struct {
	CheckedAt time.Time `json:"checked_at"`
	Available bool      `json:"available"`
}

var _ Plugin = &MemoryPlugin{}

// MemoryPlugin is an inproc plugin
type MemoryPlugin struct {
	Name        string
	Description string

	CanCheck func(*Feed) error
	Check    func(f *Feed, since time.Time) ([]Post, error)
}

// CanCheckFeed writes out things
func (mp *MemoryPlugin) CanCheckFeed(f *Feed) error {
	return mp.CanCheck(f)
}

// CheckFeed returns all posts since a given time
func (mp *MemoryPlugin) CheckFeed(f *Feed, since time.Time) ([]Post, error) {
	return mp.Check(f, since)
}

// Healthcheck always returns true for inproc plugins (is the plugin up, not site)
func (mp *MemoryPlugin) Healthcheck() (*Healthcheck, error) {
	return &Healthcheck{
		CheckedAt: time.Now(),
		Available: true,
	}, nil
}
