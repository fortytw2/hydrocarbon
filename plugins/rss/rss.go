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
)

// Reader implements hydrocarbon.Plugin for RSS feeds
type Reader struct {
	Client *http.Client
}

// Name returns the name of the plugin
func (r *Reader) Name() string {
	return "rss"
}

// Info sanitizes a URL send to this plugin, maybe by making a test request
func (r *Reader) Info(ctx context.Context, inputURL string) (title, baseURL string, err error) {
	f, err := r.getFeed(ctx, inputURL)
	if err != nil {
		return "", "", err
	}

	return f.Title, strings.TrimSpace(inputURL), nil
}

// Fetch performs the actual scraping shenanigans
func (r *Reader) Fetch(ctx context.Context, baseURL string, _ time.Time) ([]*hydrocarbon.Post, error) {
	f, err := r.getFeed(ctx, baseURL)
	if err != nil {
		return nil, err
	}

	return parseFeed(f)
}

func (r *Reader) getFeed(ctx context.Context, url string) (*Feed, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", "hydrocarbon/1.0 (+https://github.com/fortytw2/hydrocarbon)")

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpx.DrainAndClose(resp.Body)

	ct := resp.Header.Get("Content-Type")
	if !(strings.Contains(ct, "application/rss+xml") || strings.Contains(ct, "text/xml")) {
		return nil, fmt.Errorf("url has content type: %s - are you sure you have the write URL", ct)
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

		sanitized := bluemonday.UGCPolicy().Sanitize(i.Content)
		posts = append(posts, &hydrocarbon.Post{
			Author:      strings.TrimSpace(i.Author),
			Title:       strings.TrimSpace(i.Title),
			Body:        strings.TrimSpace(sanitized),
			OriginalURL: strings.TrimSpace(i.Link),
		})
	}

	return posts, nil
}
