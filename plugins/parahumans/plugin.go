package parahumans

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Puerkitobio/goquery"
	dc "github.com/fortytw2/discollect"
	"github.com/fortytw2/hydrocarbon/httpx"
	"github.com/lunny/html2md"
)

// Plugin is a plugin that can scrape parahumans
var Plugin = &dc.Plugin{
	Name: "parahumans",
	ConfigValidator: func(c *dc.Config) error {
		return nil
	},
	Routes: map[string]dc.Handler{
		`https:\/\/parahumans.wordpress.com\/(\d+)\/(\d+)\/(\d+)\/(.*)`: phPage,
	},
}

type phChapter struct {
	Author    string    `json:"author,omitempty"`
	PostedAt  time.Time `json:"posted_at,omitempty"`
	Title     string    `json:"title,omitempty"`
	Body      string    `json:"body,omitempty"`
	WordCount int       `json:"word_count,omitempty"`
}

func phPage(ctx context.Context, ho *dc.HandlerOpts, t *dc.Task) *dc.HandlerResponse {
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

	title := strings.TrimSpace(doc.Find(".entry-title").Text())

	date := strings.TrimSpace(doc.Find(".entry-date").Text())
	dateTs, err := time.Parse("January _2, 2006", date)
	if err != nil {
		return dc.ErrorResponse(err)
	}

	h := doc.Find(".entry-content")

	nextPageURL, ok := doc.Find(".nav-next").Find("a").First().Attr("href")
	if !ok {
		return dc.ErrorResponse(errors.New("no url"))
	}

	h.Find("p > a").Remove()
	h.Find(".entry-meta , #jp-post-flair , .sd-sharing").Remove()

	body, err := h.Html()
	if err != nil {
		return dc.ErrorResponse(err)
	}

	body = html2md.Convert(strings.TrimSpace(body))

	return dc.Response([]interface{}{
		&phChapter{
			Author:    "wildbow",
			PostedAt:  dateTs,
			Title:     title,
			Body:      strings.Replace(strings.TrimSpace(body), `  `, ` `, -1),
			WordCount: len(strings.Split(body, " ")),
		},
	}, &dc.Task{
		URL: nextPageURL,
	})
}
