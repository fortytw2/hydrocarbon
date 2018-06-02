package parahumans

import (
	"context"
	"errors"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/Puerkitobio/goquery"
	"github.com/fortytw2/hydrocarbon"
	dc "github.com/fortytw2/hydrocarbon/discollect"
	"github.com/fortytw2/hydrocarbon/httpx"
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

	return dc.Response([]interface{}{
		&hydrocarbon.Post{
			Author:      "wildbow",
			CreatedAt:   dateTs,
			OriginalURL: t.URL,
			Title:       title,
			Body:        html.UnescapeString(strings.Replace(strings.TrimSpace(body), `  `, ` `, -1)),
		},
	}, &dc.Task{
		URL: nextPageURL,
	})
}
