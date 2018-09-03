package discollect

import (
	"context"
	"fmt"
	"time"
)

// Resolver watches for scrapes that should be marked complete.
type Resolver struct {
	q  Queue
	ms Metastore
	er ErrorReporter

	shutdown chan chan struct{}
	ticker   *time.Ticker
}

// Start launches the resolver, which marks scrapes as complete
func (r *Resolver) Start() {
	r.ticker = time.NewTicker(pollInterval)

	for {
		select {
		case a := <-r.shutdown:
			r.ticker.Stop()
			a <- struct{}{}
			return
		case <-r.ticker.C:
			scrapes, err := r.ms.ListScrapes(context.TODO(), "RUNNING", 500, 0)
			if err != nil {
				r.er.Report(context.TODO(), nil, err)
				continue
			}

			for _, sc := range scrapes {
				ss, err := r.q.Status(context.TODO(), sc.ID)
				if err != nil {
					continue
				}

				if ss.InFlightTasks == 0 && (ss.CompletedTasks == ss.TotalTasks) {
					err = r.ms.EndScrape(context.TODO(), sc.ID, 0, ss.RetriedTasks, ss.CompletedTasks)
					if err != nil {
						continue
					}

					err = r.q.CompleteScrape(context.TODO(), sc.ID)
					if err != nil {
						// TODO(fortytw2):
						r.er.Report(context.TODO(), nil, fmt.Errorf("could not clean up redis for scrape id: %s: %s", sc.ID, err))
						continue
					}
				}
			}
		}
	}
}

// Stop gracefully stops the scheduler and blocks until its shutdown
func (r *Resolver) Stop() {
	c := make(chan struct{})
	r.shutdown <- c
	<-c
}
