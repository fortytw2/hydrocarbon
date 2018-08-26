// package redis implements a lightweight queue on top of BRPOPLPUSH
// for hydrocarbon to use
package redis

import "context"

type Queue struct {
	r *redis.Conn
}

func (q *Queue) Pop(ctx context.Context) (*QueuedTask, error) {
	return nil, nil
}

func (q *Queue) Push(ctx context.Context, tasks []*QueuedTask) error {
	return nil
}

func (q *Queue) Finish(ctx context.Context, taskID uuid.UUID) error {
	return nil
}

func (q *Queue) Error(ctx context.Context, qt *QueuedTask) error {
	return nil
}

func (q *Queue) Status(ctx context.Context, scrapeID uuid.UUID) *ScrapeStatus {
	return nil
}
