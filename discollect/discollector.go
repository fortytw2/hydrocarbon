package discollect

import (
	"context"
	"errors"
	"log"
	"sync"
)

// A Discollector ties every element of Discollect together
type Discollector struct {
	w  Writer
	r  *Registry
	l  Limiter
	ro Rotator
	q  Queue
	ms Metastore
	fs FileStore
	er ErrorReporter

	resolver *Resolver
	s        *Scheduler

	workerMu sync.RWMutex
	workers  []*Worker
}

// An OptionFn is used to pass options to a Discollector
type OptionFn func(d *Discollector) error

var defaultOpts = []OptionFn{
	WithWriter(&StdoutWriter{}),
	WithErrorReporter(&StdoutReporter{}),
	WithLimiter(&NilLimiter{}),
	WithRotator(NewDefaultRotator()),
	WithQueue(NewMemQueue()),
	WithFileStore(NewStubFS()),
}

// New returns a new Discollector
func New(opts ...OptionFn) (*Discollector, error) {
	d := &Discollector{}

	for _, o := range defaultOpts {
		err := o(d)
		if err != nil {
			return nil, err
		}
	}

	for _, o := range opts {
		err := o(d)
		if err != nil {
			return nil, err
		}
	}

	if d.r == nil {
		return nil, errors.New("no plugins registered")
	}

	d.workers = make([]*Worker, 0)

	d.s = &Scheduler{
		shutdown: make(chan chan struct{}),
		r:        d.r,
		ms:       d.ms,
		q:        d.q,
		er:       d.er,
	}

	d.resolver = &Resolver{
		ms: d.ms,
		q:  d.q,
		er: d.er,
	}

	return d, nil
}

// Start starts the scraping loops
func (d *Discollector) Start(workers int) error {
	go d.s.Start()
	go d.resolver.Start()

	d.workerMu.Lock()
	for i := workers; i > 0; i-- {
		w := NewWorker(d.r, d.ro, d.l, d.q, d.fs, d.w, d.er)
		d.workers = append(d.workers, w)
	}
	d.workerMu.Unlock()

	var wg sync.WaitGroup
	for _, w := range d.workers {
		wg.Add(1)
		go w.Start(&wg)
	}
	wg.Wait()

	return nil
}

// GetPlugin returns the plugin with the given name
func (d *Discollector) GetPlugin(name string) (*Plugin, error) {
	return d.r.Get(name)
}

// Shutdown spins down all the workers after allowing them to finish
// their current tasks
func (d *Discollector) Shutdown(ctx context.Context) {
	log.Println("stopping scheduler")
	d.s.Stop()

	log.Println("stopping scrape resolver")
	d.resolver.Stop()

	d.workerMu.Lock()
	defer d.workerMu.Unlock()

	log.Println("stopping workers")
	for _, w := range d.workers {
		w.Stop()
	}
}

// WithPlugins registers a list of plugins
func WithPlugins(p ...*Plugin) OptionFn {
	return func(d *Discollector) error {
		reg, err := NewRegistry(p)
		if err != nil {
			return err
		}

		d.r = reg

		return nil
	}
}

// WithWriter sets the Writer for the Discollector
func WithWriter(w Writer) OptionFn {
	return func(d *Discollector) error {
		d.w = w
		return nil
	}
}

// WithFileStore sets the FileStore for the Discollector
func WithFileStore(fs FileStore) OptionFn {
	return func(d *Discollector) error {
		d.fs = fs
		return nil
	}
}

// WithErrorReporter sets the ErrorReporter for the Discollector
func WithErrorReporter(er ErrorReporter) OptionFn {
	return func(d *Discollector) error {
		d.er = er
		return nil
	}
}

// WithLimiter sets the Limiter for the Discollector
func WithLimiter(l Limiter) OptionFn {
	return func(d *Discollector) error {
		d.l = l
		return nil
	}
}

// WithRotator sets the Rotator for the Discollector
func WithRotator(ro Rotator) OptionFn {
	return func(d *Discollector) error {
		d.ro = ro
		return nil
	}
}

// WithQueue sets the Queue for the Discollector
func WithQueue(q Queue) OptionFn {
	return func(d *Discollector) error {
		d.q = q
		return nil
	}
}

// WithMetastore sets the Metastore for the Discollector
func WithMetastore(ms Metastore) OptionFn {
	return func(d *Discollector) error {
		d.ms = ms
		return nil
	}
}

// ListPlugins lists all registered plugins
func (d *Discollector) ListPlugins() []string {
	var out []string
	for _, p := range d.r.plugins {
		out = append(out, p.Name)
	}

	return out
}

// StartScrape launches a new scrape
func (d *Discollector) StartScrape(ctx context.Context, pluginName string, config *Config) (string, error) {
	return "", nil
}
