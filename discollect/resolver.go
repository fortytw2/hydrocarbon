package discollect

import (
	"context"
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
	r.shutdown = make(chan chan struct{})

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
				ss := r.q.Status(context.TODO(), sc.ID)
				if ss.InFlightTasks == 0 && (ss.CompletedTasks == ss.TotalTasks) {
					r.ms.EndScrape(context.TODO(), sc.ID, 0, ss.RetriedTasks, ss.CompletedTasks)
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
