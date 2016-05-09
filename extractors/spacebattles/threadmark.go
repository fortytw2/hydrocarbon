package spacebattles

import (
	"html"
	"net/http"
	"strings"
	"time"

	"gopkg.in/inconshreveable/log15.v2"
	"github.com/Puerkitobio/goquery"
	"github.com/fortytw2/kiasu"
	"github.com/microcosm-cc/bluemonday"
)

type extractor struct{}

// NewExtractor creates a new extractor for SpaceBattles
func NewExtractor() kiasu.Extractor {
	return &extractor{}
}

func (e *extractor) Validate(f *kiasu.Feed) error {
	// URLs MUST LOOK LIKE https://forums.spacebattles.com/threads/dominion-worm-s9-taylor.340669/threadmarks
	return nil
}

func (e *extractor) Update(a *kiasu.Article) error {
	return nil
}

func (e *extractor) FindSince(f *kiasu.Feed, since time.Time) ([]kiasu.Article, error) {
	rsp, err := http.Get(f.BaseURL)
	if err != nil {
		return nil, err
	}

	log15.Info("status", "status", rsp.StatusCode)

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		return nil, err
	}

	rsp.Body.Close()

	var articles []kiasu.Article
	for _, u := range e.getThreadmarkURLs(doc, since) {
		a, err := e.getThreadmarkArticle(u)
		if err != nil {
			return nil, err
		}

		if a != nil {
			a.Content = strings.Replace(a.Content, "â", "", -1)
			a.Content = strings.Replace(a.Content, "â¦", "...", -1)

			log15.Info(a.Content)
			articles = append(articles, *a)
		}
	}

	return nil, nil
}

func (e *extractor) FindAll(f *kiasu.Feed) ([]kiasu.Article, error) {
	return e.FindSince(f, time.Time{})
}

func (e *extractor) getThreadmarkURLs(doc *goquery.Document, since time.Time) []string {
	var chapterURLs []string
	doc.Find(".threadmarkItem").EachWithBreak(func(i int, sel *goquery.Selection) bool {
		log15.Info("checking")
		ts := strings.TrimSpace(sel.Find(".DateTime").Text())
		t, err := time.Parse("Jan 2, 2006", ts)
		if err != nil {
			log15.Warn("could not parse time", "time", ts)
			return true
		}

		// respect time.Time
		if t.Before(since) {
			log15.Info("not checking", "date", ts)
			return false
		}

		url, ok := sel.Find("a").Attr("href")
		if !ok {
			log15.Warn("could not get URL for chapter")
			return false
		}

		chapterURLs = append(chapterURLs, "https://forums.spacebattles.com/"+url)
		return true
	})

	return chapterURLs
}

func (e *extractor) getThreadmarkArticle(url string) (*kiasu.Article, error) {
	split := strings.Split(url, "#post-")
	if len(split) != 2 {
		log15.Warn(url)
		// not a real article
		return nil, nil
	}

	postID := split[1]
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		return nil, err
	}

	sel := doc.Find("#post-" + postID)
	h, err := sel.Find(".messageContent").Html()
	if err != nil {
		return nil, err
	}

	p := bluemonday.UGCPolicy()
	return &kiasu.Article{
		Content: html.UnescapeString(p.Sanitize(h)),
	}, nil
}
