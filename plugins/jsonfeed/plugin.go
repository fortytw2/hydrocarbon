package jsonfeed

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/fortytw2/hydrocarbon"
	dc "github.com/fortytw2/hydrocarbon/discollect"
	"github.com/fortytw2/hydrocarbon/httpx"
	"github.com/microcosm-cc/bluemonday"
)

var rssPolicy = bluemonday.UGCPolicy().AddTargetBlankToFullyQualifiedLinks(true)

// Plugin is a plugin that can scrape rss feeds
// TODO:
// - [ ] strip images that don't matter
var Plugin = &dc.Plugin{
	Name: "jsonfeed",
	ConfigValidator: func(ho *dc.HandlerOpts) (string, error) {
		f, err := getFeed(context.TODO(), ho.Client, ho.Config.Entrypoints[0])
		if err != nil {
			return "", err
		}

		return f.Title, nil
	},
	Scheduler: dc.DefaultScheduler,
	Routes: map[string]dc.Handler{
		`(.*)`: jsonFeed,
	},
}

func jsonFeed(ctx context.Context, ho *dc.HandlerOpts, t *dc.Task) *dc.HandlerResponse {
	f, err := getFeed(ctx, ho.Client, t.URL)
	if err != nil {
		return dc.ErrorResponse(err)
	}

	posts, err := parseFeed(f)
	if err != nil {
		return dc.ErrorResponse(err)
	}

	out := make([]interface{}, len(posts))
	for i, p := range posts {
		downloaded, err := dc.DownloadImages(p.Body, ho.Client, ho.FileStore)
		if err != nil {
			return dc.ErrorResponse(err)
		}

		p.Body = downloaded

		out[i] = p
	}

	return &dc.HandlerResponse{
		Facts: out,
	}
}

func getFeed(ctx context.Context, c *http.Client, url string) (*JSONFeed, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", "hydrocarbon/1.0 (+https://github.com/fortytw2/hydrocarbon)")

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpx.DrainAndClose(resp.Body)

	var f JSONFeed
	err = json.NewDecoder(resp.Body).Decode(&f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func parseFeed(f *JSONFeed) ([]*hydrocarbon.Post, error) {
	posts := make([]*hydrocarbon.Post, 0)
	for _, i := range f.Items {
		sanitized := rssPolicy.Sanitize(i.ContentHTML)
		if sanitized == "" {
			sanitized = rssPolicy.Sanitize(i.ContentText)
		}

		if sanitized == "" {
			sanitized = rssPolicy.Sanitize(i.Summary)
		}

		var date time.Time
		if i.PublishedDate != nil {
			date = *i.PublishedDate
		}

		posts = append(posts, &hydrocarbon.Post{
			PostedAt:    date,
			Author:      strings.TrimSpace(i.Author.Name),
			Title:       strings.TrimSpace(i.Title),
			Body:        strings.TrimSpace(sanitized),
			OriginalURL: strings.TrimSpace(i.ExternalURL),
		})
	}

	return posts, nil
}
