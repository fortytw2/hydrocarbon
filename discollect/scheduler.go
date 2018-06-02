package discollect

import (
	"context"
	"time"
)

// A ScheduleStore is used to store and manage schedules of configs that need to be run
// periodically
type ScheduleStore interface {
	// ConfigToStart returns a *Config of a scrape that needs to be started
	// and the plugin to start it on
	ConfigToStart(context.Context) (string, *Config, error)

	// UpsertSchedule creates a schedule out of the given config and cron syntax
	// if it doesn't already exist
	UpsertSchedule(context.Context, *Schedule) error
}

// A Schedule is part of every plugin and defines when it needs to be run
type Schedule struct {
	Config string
	Cron   string
}

const pollInterval = 100 * time.Millisecond

// A Scheduler initiates new scrapes according to plugin-level schedules
type Scheduler struct {
	r  *Registry
	ss ScheduleStore
	ms Metastore
	q  Queue
	er ErrorReporter

	active chan chan struct{}
}

// Start launches the scheduler
func (s *Scheduler) Start() {
	for {
		select {
		case a := <-s.active:
			a <- struct{}{}
		default:
			time.Sleep(pollInterval)

			plug, conf, err := s.ss.ConfigToStart(context.TODO())
			if err != nil {
				s.er.Report(context.TODO(), nil, err)
				continue
			}

			p, err := s.r.Get(plug)
			if err != nil {
				s.er.Report(context.TODO(), nil, err)
				continue
			}

			err = launchScrape(context.TODO(), p, conf, s.q, s.ms)
			if err != nil {
				s.er.Report(context.TODO(), nil, err)
			}
		}
	}
}

// Stop gracefully stops the scheduler and blocks until its shutdown
func (s *Scheduler) Stop() {
	c := make(chan struct{})
	s.active <- c
	<-c
}
