package discollect

import (
	"context"
	"fmt"
	"time"
)

const pollInterval = 30 * time.Second
const scrapeLimit = 25
const forwardScrapeLimit = 5

// A Scheduler initiates new scrapes according to plugin-level schedules
type Scheduler struct {
	r  *Registry
	ms Metastore
	q  Queue
	er ErrorReporter

	ticker   *time.Ticker
	shutdown chan chan struct{}
}

// Start launches the scheduler
func (s *Scheduler) Start() {
	s.shutdown = make(chan chan struct{})

	s.ticker = time.NewTicker(pollInterval)

	for {
		select {
		case a := <-s.shutdown:
			s.ticker.Stop()
			a <- struct{}{}
			return
		case <-s.ticker.C:
			scrapes, err := s.ms.StartScrapes(context.TODO(), scrapeLimit)
			if err != nil {
				s.er.Report(context.TODO(), nil, err)
				continue
			}

			for _, sc := range scrapes {
				p, err := s.r.Get(sc.Plugin)
				if err != nil {
					s.er.Report(context.TODO(), nil, err)
					continue
				}

				err = launchScrape(context.TODO(), sc.ID, p, sc.Config, s.q, s.ms)
				if err != nil {
					s.er.Report(context.TODO(), nil, err)
				}
			}

			// forward scheduler action
			srs, err := s.ms.FindMissingSchedules(context.TODO(), forwardScrapeLimit)
			if err != nil {
				s.er.Report(context.TODO(), nil, err)
				continue
			}

			for _, sr := range srs {
				p, err := s.r.Get(sr.Plugin)
				if err != nil {
					s.er.Report(context.TODO(), nil, fmt.Errorf("forward-scheduler: cannot find plugin: %s", err))
					continue
				}

				ss, err := p.Scheduler(sr)
				if err != nil {
					s.er.Report(context.TODO(), nil, err)
					continue
				}

				err = s.ms.InsertSchedule(context.TODO(), sr, ss)
				if err != nil {
					s.er.Report(context.TODO(), nil, err)
				}
			}
		}
	}
}

// Stop gracefully stops the scheduler and blocks until its shutdown
func (s *Scheduler) Stop() {
	c := make(chan struct{})
	s.shutdown <- c
	<-c
}
