package xenforo

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/fortytw2/hydrocarbon"
	dc "github.com/fortytw2/hydrocarbon/discollect"
	"github.com/fortytw2/hydrocarbon/httpx"

	"github.com/Puerkitobio/goquery"
	"github.com/microcosm-cc/bluemonday"
)

var firstPost = "first"

// Plugin is a plugin that can scrape xenforo threads by threadmark
var Plugin = &dc.Plugin{
	Name: "xenforo_threadmarks",
	ConfigValidator: func(ho *dc.HandlerOpts) (string, error) {
		if !strings.Contains(ho.Config.Entrypoints[0], "/threadmarks") {
			return "", errors.New("entrypoint for xenforo_threadmarks MUST end in /threadmarks")
		}

		page, err := threadmarksPage(ho.Client, ho.Config.Entrypoints[0])
		if err != nil {
			return "", err
		}

		title := strings.TrimSpace(page.Find(`.titleBar > h1:nth-child(1)`).Text())
		title = strings.TrimPrefix(title, "Threadmarks for: ")

		return title, nil
	},
	Scheduler: dc.DefaultScheduler,
	Routes: map[string]dc.Handler{
		`(.*)/threadmarks`:    threadmarksHandler,
		`(.*)/page-(\d+)(.*)`: postHandler,
	},
}

func threadmarksHandler(ctx context.Context, ho *dc.HandlerOpts, t *dc.Task) *dc.HandlerResponse {

	return nil
}

func postHandler(ctx context.Context, ho *dc.HandlerOpts, t *dc.Task) *dc.HandlerResponse {

	return nil
}

func threadmarksPage(c *http.Client, url string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer httpx.DrainAndClose(resp.Body)

	return doc, nil
}

func getThreadmarkURLs(doc *goquery.Document, since time.Time) []string {
	var chapterURLs []string
	doc.Find(".threadmarkItem").EachWithBreak(func(i int, sel *goquery.Selection) bool {
		ts := strings.TrimSpace(sel.Find(".DateTime").Text())
		tss := strings.Split(ts, " at ")

		t, err := time.Parse("Jan 2, 2006", tss[0])
		if err != nil {
			return true
		}

		// respect time.Time
		if t.Before(since) {
			return false
		}

		url, ok := sel.Find("a").Attr("href")
		if !ok {

			return false
		}

		chapterURLs = append(chapterURLs, "https://forums.spacebattles.com/"+url)
		return true
	})

	return chapterURLs
}

func getThreadmarkPost(c *http.Client, url string) (*hydrocarbon.Post, error) {
	split := strings.Split(url, "#post-")

	var postID string
	if len(split) != 2 {
		postID = firstPost
	} else {
		postID = split[1]
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpx.DrainAndClose(resp.Body)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var sel *goquery.Selection
	if postID == firstPost {
		sel = doc.Find(".hasThreadmark").First()
	} else {
		sel = doc.Find("#post-" + postID)
	}

	h, err := sel.Find(".messageContent").Html()
	if err != nil {
		return nil, err
	}

	var postTime time.Time
	t, ok := sel.Find(".DateTime").Attr("title")
	if ok {
		postTime, err = time.Parse("Jan 2, 2006 at 3:04 PM", t)
		if err != nil {
			return nil, err
		}
	}

	title := strings.Replace(sel.Find(".threadmarker > .label").Text(), "Threadmark:", "", -1)
	title = strings.TrimSpace(title)

	p := bluemonday.UGCPolicy()
	html := p.Sanitize(h)

	return &hydrocarbon.Post{
		OriginalURL: url,
		PostedAt:    postTime,
		Title:       title,
		Body:        html,
	}, nil
}
