package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/httpx"
	"github.com/microcosm-cc/bluemonday"

	dc "github.com/fortytw2/hydrocarbon/discollect"
)

// Sun, 15 Jul 2018 11:48:43 -0700
const rssTime = "Mon, 02 Jan 2006 15:04:05 -0700"
const rssAltTime = "Mon, 02 Jan 2006 15:04:05 MST"

var rssPolicy = bluemonday.UGCPolicy().AddTargetBlankToFullyQualifiedLinks(true)

// Plugin is a plugin that can scrape rss feeds
// TODO:
// - [ ] strip images that don't matter
// - [ ] download images
var Plugin = &dc.Plugin{
	Name: "rss",
	ConfigValidator: func(c *dc.Config) error {
		return nil
	},
	Routes: map[string]dc.Handler{
		`(.*)`: rssFeed,
	},
}

func rssFeed(ctx context.Context, ho *dc.HandlerOpts, t *dc.Task) *dc.HandlerResponse {
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
		out[i] = p
	}

	return &dc.HandlerResponse{
		Facts: out,
	}
}

func getFeed(ctx context.Context, c *http.Client, url string) (*Feed, error) {
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

	ct := resp.Header.Get("Content-Type")
	if !(strings.Contains(ct, "application/rss+xml") || strings.Contains(ct, "text/xml")) || strings.Contains(ct, "application/xml") {
		return nil, fmt.Errorf("url has content type: %s - are you sure you have the correct URL", ct)
	}

	var f Feed
	err = xml.NewDecoder(resp.Body).Decode(&f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func parseFeed(f *Feed) ([]*hydrocarbon.Post, error) {
	posts := make([]*hydrocarbon.Post, 0)
	for _, i := range f.Items {

		var pubDate time.Time
		if f.PubDate != "" {
			var err error
			pubDate, err = time.Parse(rssTime, f.PubDate)
			if err != nil {
				pubDate, err = time.Parse(rssAltTime, f.PubDate)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		if pubDate.IsZero() {
			pubDate = time.Now()
		}

		sanitized := rssPolicy.Sanitize(i.Content)
		if sanitized == "" {
			sanitized = rssPolicy.Sanitize(i.Description)
		}

		posts = append(posts, &hydrocarbon.Post{
			PostedAt:    pubDate,
			Author:      strings.TrimSpace(i.Author),
			Title:       strings.TrimSpace(i.Title),
			Body:        strings.TrimSpace(sanitized),
			OriginalURL: strings.TrimSpace(i.Link),
		})
	}

	return posts, nil
}
