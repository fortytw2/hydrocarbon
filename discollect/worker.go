package discollect

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// A Worker is a single-threaded worker that pulls a single task from the queue at a time
// and process it to completion
type Worker struct {
	r  *Registry
	ro Rotator
	l  Limiter
	q  Queue
	w  Writer
	fs FileStore
	er ErrorReporter

	shutdown chan chan struct{}
}

// NewWorker provisions a new worker
func NewWorker(r *Registry, ro Rotator, l Limiter, q Queue, fs FileStore, w Writer, er ErrorReporter) *Worker {
	return &Worker{
		r:        r,
		ro:       ro,
		l:        l,
		q:        q,
		fs:       fs,
		w:        w,
		er:       er,
		shutdown: make(chan chan struct{}),
	}
}

// Start launches the worker
func (w *Worker) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case s := <-w.shutdown:
			s <- struct{}{}
			return
		default:
			qt, err := w.q.Pop(context.TODO())
			if err != nil {
				w.er.Report(context.TODO(), nil, fmt.Errorf("discollect: worker-pop: %s", err))
				continue
			}

			if qt == nil {
				time.Sleep(time.Second * 1)
				continue
			}

			timeout := defaultTimeout
			if qt.Task.Timeout != 0 {
				timeout = qt.Task.Timeout
			}

			// set config timeout on all worker actions on this task
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err = w.processTask(ctx, qt)
			if err != nil {
				w.er.Report(ctx, nil, fmt.Errorf("discollect: worker-process-task: %s", err))
				// retry task
				w.q.Error(ctx, qt)
				cancel()
				continue
			}

			// callback that we've finished the task
			err = w.q.Finish(ctx, qt)
			if err != nil {
				w.er.Report(ctx, nil, err)
			}

			cancel()
		}
	}
}

// Stop initiates stop and then blocks until shutdown is complete
func (w *Worker) Stop() {
	c := make(chan struct{})
	w.shutdown <- c
	<-c
}

// processTask executes one task
// Safe for concurrent use.
func (w *Worker) processTask(ctx context.Context, q *QueuedTask) error {
	handler, params, err := w.r.HandlerFor(q.Plugin, q.Task.URL)
	if err != nil {
		return err
	}

	plugin, err := w.r.Get(q.Plugin)
	if err != nil {
		return err
	}

	// if this rate limit blocks too long and the context cancels we can just
	// return error and the task will be retried later
	res, err := w.l.Reserve(plugin.RateLimit, q.Task.URL, q.ScrapeID)
	if err != nil {
		return err
	}

	if res.OK() {
		if res.Delay() < time.Second*10 {
			time.Sleep(res.Delay())
		} else {
			res.Cancel()
			return ErrRateLimitExceeded
		}
	} else {
		return ErrRateLimitExceeded
	}

	client, err := w.ro.Get(q.Config)
	if err != nil {
		return err
	}

	resp := handler(ctx, &HandlerOpts{
		Config:      q.Config,
		FileStore:   w.fs,
		RouteParams: params,
		Client:      client,
	}, q.Task)

	// report errors
	for _, err := range resp.Errors {
		w.er.Report(ctx, &ReporterOpts{
			ScrapeID: q.ScrapeID,
			Plugin:   q.Plugin,
			URL:      q.Task.URL,
		}, err)
	}

	// submit tasks back to the queue
	qt := make([]*QueuedTask, 0)
	for _, t := range resp.Tasks {
		if t == nil {
			continue
		}

		qt = append(qt, &QueuedTask{
			ScrapeID: q.ScrapeID,
			Plugin:   q.Plugin,
			Config:   q.Config,
			QueuedAt: time.Now().In(time.UTC),
			TaskID:   uuid.New(),
			Task:     t,
		})
	}

	if len(qt) > 0 {
		err = w.q.Push(ctx, qt)
		if err != nil {
			return err
		}
	}

	// write facts
	for _, f := range resp.Facts {
		err = w.w.Write(ctx, q.ScrapeID, f)
		if err != nil {
			return err
		}
	}

	return nil
}
