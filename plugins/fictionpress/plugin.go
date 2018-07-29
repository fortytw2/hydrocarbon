package fictionpress

import (
	"context"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Puerkitobio/goquery"
	"github.com/fortytw2/hydrocarbon"
	dc "github.com/fortytw2/hydrocarbon/discollect"

	"github.com/fortytw2/hydrocarbon/httpx"
)

// Plugin is a plugin that can scrape fictionpress
var Plugin = &dc.Plugin{
	Name: "fictionpress",
	ConfigValidator: func(ho *dc.HandlerOpts) (string, error) {
		for _, e := range ho.Config.Entrypoints {
			if !strings.Contains(e, "fictionpress.com") && !strings.Contains(e, "fanfiction.net") {
				return "", errors.New("fictionpress plugin only works for fictionpress and fanfiction.net")
			}
		}

		return getTitle(ho)
	},
	Routes: map[string]dc.Handler{
		`https:\/\/www.(fictionpress.com|fanfiction.net)\/s\/(.*)\/(\d+)(.*)`: storyPage,
	},
}

func getTitle(ho *dc.HandlerOpts) (string, error) {
	resp, err := ho.Client.Get(ho.Config.Entrypoints[0])
	if err != nil {
		return "", err
	}
	defer httpx.DrainAndClose(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("non-200 response")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(doc.Find("#profile_top > b").First().Text()), nil
}

func storyPage(ctx context.Context, ho *dc.HandlerOpts, t *dc.Task) *dc.HandlerResponse {
	parsedURL, err := url.Parse(t.URL)
	if err != nil {
		return dc.ErrorResponse(err)
	}

	resp, err := ho.Client.Get(t.URL)
	if err != nil {
		return dc.ErrorResponse(err)
	}
	defer httpx.DrainAndClose(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return dc.ErrorResponse(errors.New("did not get 200"))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return dc.ErrorResponse(err)
	}

	body, err := doc.Find(`#storytext`).Html()
	if err != nil {
		return dc.ErrorResponse(err)
	}

	title := doc.Find(`#chap_select > option[selected]`).First().Text()
	titleSplit := strings.Split(title, ". ")
	if len(titleSplit) != 2 {
		return dc.ErrorResponse(errors.New("could not find title or number"))
	}

	chapter := titleSplit[0]
	chapterTitle := titleSplit[1]
	day, err := strconv.Atoi(chapter)
	if err != nil {
		return dc.ErrorResponse(err)
	}

	c := &hydrocarbon.Post{
		// there is no posted at date for either site, so make a fake date using
		// the year to maintain ordering
		PostedAt: time.Date(day, 01, 01, 0, 0, 0, 0, time.UTC),
		Title:    chapterTitle,
		Author:   strings.TrimSpace(doc.Find(`#profile_top .xcontrast_txt+ a.xcontrast_txt`).Text()),
		Body:     html.UnescapeString(strings.TrimSpace(body)),
	}

	// find all chapters if this is the first one
	var tasks []*dc.Task
	// only for the first task
	if ho.RouteParams[3] == "1" {
		doc.Find(`#chap_select`).First().Find(`option`).Each(func(i int, sel *goquery.Selection) {
			val, exists := sel.Attr("value")
			if !exists || val == "1" {
				return
			}

			tasks = append(tasks, &dc.Task{
				URL: fmt.Sprintf("https://%s/s/%s/%s", parsedURL.Host, ho.RouteParams[2], val),
			})
		})
	}

	return &dc.HandlerResponse{
		Facts: []interface{}{
			c,
		},
		Tasks: tasks,
	}
}
