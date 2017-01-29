package hydrocarbon

import (
	"time"

	"github.com/fortytw2/hydrocarbon/internal/log"
	multierror "github.com/hashicorp/go-multierror"
)

// ScrapeLoop starts the background scraper
func ScrapeLoop(l log.Logger, s *Store, plugins map[string]Instantiator) {
	l.Log("msg", "starting scrape loop")
	c := DefaultClient()

	for {
		feeds, err := s.Feeds.GetFeedsToUpdate(1)
		if err != nil {
			l.Log("msg", "could not get feeds to refresh", "error", err)
			continue
		}

		for _, f := range feeds {
			l.Log("msg", "updating feed", "feed", f.Name, "id", f.ID, "plugin", f.Plugin)
			err := Scrape(plugins[f.Plugin], f, s, c)
			if err != nil {
				l.Log("msg", "could not update feed", "error", err)
			}
			err = s.Feeds.SetRefreshedAt(f.ID)
			if err != nil {
				// CRITICAL
				panic(err)
			} else {
				l.Log("msg", "successfully updated feed", "feed", f.Name, "id", f.ID, "plugin", f.Plugin)
			}
		}
		time.Sleep(30 * time.Second)
	}
}

// Scrape does a complete update of a single feed
func Scrape(ist Instantiator, f Feed, s *Store, c Client) error {
	plug, err := ist()
	if err != nil {
		return err
	}

	var t time.Time
	if !f.LastRefreshedAt.Valid {
		t = time.Time{}
	} else {
		t = f.LastRefreshedAt.Time
	}

	err = plug.Validate(c, Config{InitialURL: f.InitialURL, Since: t})
	if err != nil {
		return err
	}

	posts, err := plug.Run(c, Config{InitialURL: f.InitialURL, Since: t})
	if err != nil {
		return err
	}

	var e error
	for _, p := range posts {
		p.FeedID = f.ID
		_, err := s.Posts.CreatePost(&p)
		if err != nil {
			e = multierror.Append(e, err)
		}
	}

	return e
}
