package discollect

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const defaultScheduleHorizon = time.Hour * 4

// A ScheduleRequest is used to ask for future schedules
type ScheduleRequest struct {
	Plugin        string
	FeedID        uuid.UUID
	LatestScrapes []*Scrape
	LatestDatums  interface{}
}

// A ScrapeSchedule adds to the future
type ScrapeSchedule struct {
	Config           *Config
	ScheduledStartAt time.Time
}

// DefaultScheduler uses a simple heuristic to predict when to next scrape.
func DefaultScheduler(sr *ScheduleRequest) ([]*ScrapeSchedule, error) {
	if len(sr.LatestScrapes) == 0 {
		return nil, errors.New("discollect: cannot schedule a scrape without an initial scrape")
	}

	base := time.Now()
	conf := sr.LatestScrapes[0].Config

	var ss []*ScrapeSchedule
	for x := defaultScheduleHorizon; x > 0; x -= 30 * time.Minute {
		ss = append(ss, &ScrapeSchedule{
			ScheduledStartAt: base.Add(x),
			Config:           conf,
		})
	}
	return ss, nil
}

// NeverSchedule simple never schedules another scrape
func NeverSchedule(sr *ScheduleRequest) ([]*ScrapeSchedule, error) {
	if len(sr.LatestScrapes) == 0 {
		return nil, errors.New("discollect: cannot schedule a scrape without an initial scrape")
	}

	base := time.Now()
	conf := sr.LatestScrapes[0].Config

	return []*ScrapeSchedule{
		{
			ScheduledStartAt: base.Add(time.Hour * 24 * 365),
			Config:           conf,
		},
	}, nil
}
