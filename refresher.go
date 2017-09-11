package hydrocarbon

import (
	"context"
	"time"
)

// A RefreshController is used to fetch/update feeds from a persistent store
type RefreshController interface {
	GetFeedsToRefresh(ctx context.Context, num int) ([]*Feed, error)
	UpdateFeedFromRefresh(ctx context.Context, feedID string, posts []*Post) error
}

// a Refresher is used to poll and update feeds in the background
type Refresher struct {
	rc RefreshController
	pl *PluginList
	er ErrorReporter
}

// NewRefresher constructs a new Refresher
func NewRefresher(rc RefreshController, pl *PluginList, er ErrorReporter) *Refresher {
	return &Refresher{
		rc: rc,
		pl: pl,
		er: er,
	}
}

// Refresh is the background poller
func (rf *Refresher) Refresh(ctx context.Context) {
	t := time.NewTicker(time.Second * 5)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			rf.updateFeeds(ctx)
		}
	}
}

func (rf *Refresher) updateFeeds(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	feeds, err := rf.rc.GetFeedsToRefresh(ctx, 10)
	if err != nil {
		rf.er.Report(ctx, err)
		return
	}

	for _, f := range feeds {
		plug, err := rf.pl.ByName(f.Plugin)
		if err != nil {
			rf.er.Report(ctx, err)
			continue
		}

		posts, err := plug.Fetch(ctx, f.BaseURL, f.UpdatedAt)
		if err != nil {
			rf.er.Report(ctx, err)
			continue
		}

		err = rf.rc.UpdateFeedFromRefresh(ctx, f.ID, posts)
		if err != nil {
			rf.er.Report(ctx, err)
		}
	}
}
