// package redis implements a lightweight queue on top of BRPOPLPUSH
// for hydrocarbon to use
package redis

import (
	"context"

	"github.com/fortytw2/hydrocarbon/discollect"
	"github.com/google/uuid"
)

type Queue struct {
	// r *redis.
}

func (q *Queue) Pop(ctx context.Context) (*discollect.QueuedTask, error) {
	return nil, nil
}

func (q *Queue) Push(ctx context.Context, tasks []*discollect.QueuedTask) error {
	return nil
}

func (q *Queue) Finish(ctx context.Context, taskID uuid.UUID) error {
	return nil
}

func (q *Queue) Error(ctx context.Context, qt *discollect.QueuedTask) error {
	return nil
}

func (q *Queue) Status(ctx context.Context, scrapeID uuid.UUID) *discollect.ScrapeStatus {
	return nil
}
