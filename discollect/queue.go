package discollect

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// A Queue is used to submit and retrieve individual tasks
type Queue interface {
	Pop(ctx context.Context) (*QueuedTask, error)
	Push(ctx context.Context, tasks []*QueuedTask) error

	Finish(ctx context.Context, qt *QueuedTask) error
	Error(ctx context.Context, qt *QueuedTask) error

	Status(ctx context.Context, scrapeID uuid.UUID) (*ScrapeStatus, error)

	CompleteScrape(ctx context.Context, scrapeID uuid.UUID) error
}

var ErrCompletedScrape = errors.New("completed scrape")

// A QueuedTask is the struct for a task that goes on the Queue
type QueuedTask struct {
	// set by the TaskQueue
	TaskID   uuid.UUID `json:"task_id"`
	ScrapeID uuid.UUID `json:"scrape_id"`

	QueuedAt time.Time `json:"queued_at"`
	Config   *Config   `json:"config"`
	Plugin   string    `json:"plugin"`
	Retries  int       `json:"retries"`

	Task *Task `json:"task"`
}

// A Task generally maps to a single HTTP request, but sometimes more than one
// may be made
type Task struct {
	URL string `json:"url"`
	// Extra can be used to send information from a parent task to its children
	Extra map[string]json.RawMessage `json:"extra,omitempty"`
	// Timeout is the timeout a single task should have attached to it
	// defaults to 15s
	Timeout time.Duration
}

// ScrapeStatus is returned from a Queue with information about a specific scrape
type ScrapeStatus struct {
	TotalTasks     int `json:"total_tasks,omitempty"`
	InFlightTasks  int `json:"in_flight_tasks,omitempty"`
	CompletedTasks int `json:"completed_tasks,omitempty"`
	RetriedTasks   int `json:"retried_tasks,omitempty"`
}

// NewMemQueue makes a new purely in-memory queue
func NewMemQueue() *MemQueue {
	return &MemQueue{
		state: make(map[uuid.UUID]*ScrapeStatus),
		q:     make(map[uuid.UUID]chan *QueuedTask),
	}
}

// A MemQueue is a super simple Queue backed by an array and a mutex
type MemQueue struct {
	mu sync.Mutex

	state map[uuid.UUID]*ScrapeStatus
	q     map[uuid.UUID]chan *QueuedTask
}

func (mq *MemQueue) reset() {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	mq.state = make(map[uuid.UUID]*ScrapeStatus)
	mq.q = make(map[uuid.UUID]chan *QueuedTask)
}

// Pop pops a single task off the left side of the array
func (mq *MemQueue) Pop(ctx context.Context) (*QueuedTask, error) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if len(mq.q) == 0 {
		return nil, nil
	}

	for _, q := range mq.q {
		select {
		case task := <-q:
			if task != nil {
				mq.state[task.ScrapeID].InFlightTasks += 1
				return task, nil
			} else {
				return nil, nil
			}
		default:
			continue
		}
	}

	return nil, nil
}

// Push appends tasks to the right side of the array
func (mq *MemQueue) Push(ctx context.Context, tasks []*QueuedTask) error {
	for _, t := range tasks {
		if t == nil {
			continue
		}

		mq.mu.Lock()
		if mq.state[t.ScrapeID] == nil {
			mq.state[t.ScrapeID] = &ScrapeStatus{}
		}

		if t.Retries == 0 {
			mq.state[t.ScrapeID].TotalTasks += 1
		} else {
			mq.state[t.ScrapeID].RetriedTasks += 1
		}

		var c chan *QueuedTask
		if mq.q[t.ScrapeID] == nil {
			c = make(chan *QueuedTask, 64)
			mq.q[t.ScrapeID] = c
		} else {
			c = mq.q[t.ScrapeID]
		}
		mq.mu.Unlock()

		c <- t
	}

	return nil
}

func (mq *MemQueue) Error(ctx context.Context, qt *QueuedTask) error {
	mq.mu.Lock()
	mq.state[qt.ScrapeID].InFlightTasks -= 1
	mq.state[qt.ScrapeID].RetriedTasks += 1

	writeTo := mq.q[qt.ScrapeID]
	mq.mu.Unlock()

	writeTo <- qt

	return nil
}

// Finish is a no-op for the MemQueue
func (mq *MemQueue) Finish(ctx context.Context, qt *QueuedTask) error {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	mq.state[qt.ScrapeID].InFlightTasks -= 1
	mq.state[qt.ScrapeID].CompletedTasks += 1

	return nil
}

// Status returns the status for a given scrape
func (mq *MemQueue) Status(ctx context.Context, scrapeID uuid.UUID) (*ScrapeStatus, error) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	ss := mq.state[scrapeID]
	cop := *ss

	return &cop, nil
}

func (mq *MemQueue) CompleteScrape(ctx context.Context, scrapeID uuid.UUID) error {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	delete(mq.state, scrapeID)
	delete(mq.q, scrapeID)

	return nil
}
