// package redis implements a lightweight queue on top of RPOPLPUSH
// for hydrocarbon to use
package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fortytw2/hydrocarbon/discollect"
	"github.com/garyburd/redigo/redis"
	"github.com/google/uuid"
)

const activeScrapeIDsKey = `active_scrape_ids`

func scrapeTasksKey(scrapeID uuid.UUID) string {
	return fmt.Sprintf("%s_tasks", scrapeID)
}

func scrapeInflightTasksKey(scrapeID uuid.UUID) string {
	return fmt.Sprintf("%s_inflight_tasks", scrapeID)
}

func scrapeTotalCounterKey(scrapeID uuid.UUID) string {
	return fmt.Sprintf("%s_total", scrapeID)
}

func scrapeInflightCounterKey(scrapeID uuid.UUID) string {
	return fmt.Sprintf("%s_inflight", scrapeID)
}

func scrapeRetriesCounterKey(scrapeID uuid.UUID) string {
	return fmt.Sprintf("%s_retries", scrapeID)
}

func scrapeCompletedCounterKey(scrapeID uuid.UUID) string {
	return fmt.Sprintf("%c_completed", scrapeID)
}

// Queue implements discollect.Queue using a redis reliable queue
type Queue struct {
	r *redis.Pool
}

// NewQueue instantiates a queue, checks redis, and returns
func NewQueue(redisAddr string, redisDBIndex int) (*Queue, error) {
	pool := &redis.Pool{
		// Other pool configuration not shown in this example.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisAddr)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("SELECT", redisDBIndex); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}

	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	if err != nil {
		return nil, err
	}

	return &Queue{
		r: pool,
	}, nil
}

// SRANDMEMBER active_scrape_ids
// RPOPLPUSH from scrapeid_tasks to scrapeid_inflight_tasks
// INCR scrapeid_inflight
func (q *Queue) Pop(ctx context.Context) (*discollect.QueuedTask, error) {
	conn := q.r.Get()
	defer conn.Close()

	strScrapeID, err := redis.String(conn.Do("SRANDMEMBER", activeScrapeIDsKey))
	if err != nil {
		// nothing in the queue at all
		if err == redis.ErrNil {
			return nil, nil
		}

		return nil, err
	}

	scrapeID, err := uuid.Parse(strScrapeID)
	if err != nil {
		return nil, err
	}

	task, err := redis.Bytes(conn.Do("RPOPLPUSH", scrapeTasksKey(scrapeID), scrapeInflightTasksKey(scrapeID)))
	if err != nil {
		// no task in the queue for that specific scrape
		if err == redis.ErrNil {
			return nil, nil
		}

		return nil, err
	}

	var qt discollect.QueuedTask
	err = json.Unmarshal(task, &qt)
	if err != nil {
		return nil, err
	}

	_, err = redis.Int(conn.Do("INCR", scrapeInflightCounterKey(scrapeID)))
	if err != nil {
		return nil, err
	}

	return &qt, nil
}

// Push adds a slice of tasks onto the queue
// SADD scrapeid to active_scrape_ids
// INCR scrapeid_total
// LPUSH onto 'scrapeid_tasks'
func (q *Queue) Push(ctx context.Context, tasks []*discollect.QueuedTask) error {
	if len(tasks) == 0 {
		return nil
	}

	conn := q.r.Get()
	defer conn.Close()

	scrapeID := tasks[0].ScrapeID

	_, err := redis.Int(conn.Do("SADD", activeScrapeIDsKey, scrapeID))
	if err != nil {
		return err
	}

	_, err = redis.Bool(conn.Do("INCRBY", scrapeTotalCounterKey(scrapeID), len(tasks)))
	if err != nil {
		return err
	}

	lpushSet := make([]interface{}, len(tasks)+1)
	lpushSet[0] = scrapeTasksKey(scrapeID)
	for i, t := range tasks {
		buf, err := json.Marshal(t)
		if err != nil {
			return err
		}

		lpushSet[i+1] = buf
	}

	_, err = redis.Int(conn.Do("LPUSH", lpushSet...))
	return err
}

// Finish marks a task as fully complete
// INCR scrapeid_complete
// LREM from scrapeid_inflight_tasks
func (q *Queue) Finish(ctx context.Context, task *discollect.QueuedTask) error {
	conn := q.r.Get()
	defer conn.Close()

	_, err := redis.Int(conn.Do("INCR", scrapeCompletedCounterKey(task.ScrapeID)))
	if err != nil {
		return err
	}

	_, err = redis.Int(conn.Do("DECR", scrapeInflightCounterKey(task.ScrapeID)))
	if err != nil {
		return err
	}

	buf, err := json.Marshal(task)
	if err != nil {
		return err
	}

	_, err = redis.Int(conn.Do("LREM", scrapeInflightTasksKey(task.ScrapeID), "0", buf))
	return err
}

// Error retries a given task
// INCR retries_counter
// LREM inflight-tasks
// DECR inflight_counter
// LPUSH tasks
func (q *Queue) Error(ctx context.Context, task *discollect.QueuedTask) error {
	conn := q.r.Get()
	defer conn.Close()

	_, err := redis.Int(conn.Do("INCR", scrapeRetriesCounterKey(task.ScrapeID)))
	if err != nil {
		return err
	}

	buf, err := json.Marshal(task)
	if err != nil {
		return err
	}

	_, err = redis.Int(conn.Do("LREM", scrapeInflightTasksKey(task.ScrapeID), "0", buf))
	if err != nil {
		return err
	}

	_, err = redis.Int(conn.Do("DECR", scrapeInflightCounterKey(task.ScrapeID)))
	if err != nil {
		return err
	}

	_, err = redis.Int(conn.Do("LPUSH", scrapeTasksKey(task.ScrapeID), buf))
	return err
}

// Status returns the status of a given scrape
func (q *Queue) Status(ctx context.Context, scrapeID uuid.UUID) (*discollect.ScrapeStatus, error) {
	conn := q.r.Get()
	defer conn.Close()

	vals, err := redis.Ints(conn.Do("MGET",
		scrapeTotalCounterKey(scrapeID),
		scrapeCompletedCounterKey(scrapeID),
		scrapeRetriesCounterKey(scrapeID),
		scrapeInflightCounterKey(scrapeID)))
	if err != nil {
		return nil, err
	}

	if len(vals) != 4 {
		return nil, errors.New("could not get scrape status")
	}

	return &discollect.ScrapeStatus{
		TotalTasks:     vals[0],
		CompletedTasks: vals[1],
		RetriedTasks:   vals[2],
		InFlightTasks:  vals[3],
	}, nil
}

// DELETE scrapeid_tasks
// DELETE scrapeid_inflight_tasks

// DELETE scrapeid_total
// DELETE scrapeid_complete
// DELETE scrapeid_retries
// DELETE scrapeid_inflight

// DELETE scrapeid FROM active_scrape_ids
func (q *Queue) CompleteScrape(ctx context.Context, scrapeID uuid.UUID) error {
	conn := q.r.Get()
	defer conn.Close()

	keys := []string{
		scrapeTasksKey(scrapeID),
		scrapeInflightTasksKey(scrapeID),
		scrapeTotalCounterKey(scrapeID),
		scrapeCompletedCounterKey(scrapeID),
		scrapeRetriesCounterKey(scrapeID),
		scrapeInflightCounterKey(scrapeID),
	}

	for _, k := range keys {
		_, err := redis.Bool(conn.Do("DEL", k))
		if err != nil {
			return err
		}
	}

	_, err := redis.Int(conn.Do("SREM", activeScrapeIDsKey, scrapeID))
	return err
}

// ResetAll runs FLUSHALL
func (q *Queue) ResetAll() error {
	conn := q.r.Get()
	defer conn.Close()

	_, err := conn.Do("FLUSHALL")
	return err
}
