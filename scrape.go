package hydrocarbon

import (
	"time"

	"github.com/fortytw2/hydrocarbon/internal/log"
	multierror "github.com/hashicorp/go-multierror"
	uuid "github.com/satori/go.uuid"
)

// ScrapeLoop starts the background scraper
func ScrapeLoop(l log.Logger, fs FeedStore, ps PostStore, plugins map[string]Instantiator) {
	l.Log("msg", "starting scrape loop")
	c := DefaultClient()

	for {
		feeds, err := fs.GetFeeds(&Pagination{
			Page:     0,
			PageSize: 1000,
		})
		if err != nil {
			l.Log("msg", "could not get feeds to refresh", "error", err)
			continue
		}

		for _, f := range feeds {
			l.Log("msg", "updating feed", "feed", f.Name, "plugin", f.Plugin)
			err := Scrape(plugins[f.Plugin], f, ps, c)
			if err != nil {
				l.Log("msg", "could not update feed", "error", err)
			}
		}
		time.Sleep(5 * time.Minute)
	}
}

// Scrape does a complete update of a single feed
func Scrape(ist Instantiator, f Feed, ps PostStore, c Client) error {
	plug, err := ist()
	if err != nil {
		return err
	}

	var t time.Time
	if !f.RefreshedAt.Valid {
		t = time.Time{}
	} else {
		t = f.RefreshedAt.Time
	}

	err = plug.Validate(c, Config{InitialURL: f.InitialURL, Since: t})
	if err != nil {
		return err
	}

	println("starting run")
	posts, err := plug.Run(c, Config{InitialURL: f.InitialURL, Since: t})
	if err != nil {
		return err
	}
	println("ending run")

	var e error
	for _, p := range posts {
		p.ID = uuid.NewV4().String()
		p.FeedID = f.ID
		_, err := ps.CreatePost(&p)
		if err != nil {
			e = multierror.Append(e, err)
		}
	}

	return e
}
