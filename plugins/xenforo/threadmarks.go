package xenforo

import (
	"net/http"
	"strings"
	"time"

	"github.com/Puerkitobio/goquery"
	"github.com/fortytw2/kiasu"
	"github.com/jaytaylor/html2text"
)

type Extractor struct{}

// NewExtractor creates a new Extractor for SpaceBattles
func NewExtractor() *Extractor {
	return &Extractor{}
}

func (e *Extractor) Validate(f *kiasu.Feed) error {
	// URLs MUST LOOK LIKE https://forums.spacebattles.com/threads/dominion-worm-s9-taylor.340669/threadmarks
	return nil
}

func (e *Extractor) Update(a *kiasu.Post) error {
	return nil
}

func (e *Extractor) FindSince(f *kiasu.Feed, since time.Time) ([]kiasu.Post, error) {
	rsp, err := http.Get(f.URL)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		return nil, err
	}

	rsp.Body.Close()

	var articles []kiasu.Post
	for _, u := range e.getThreadmarkURLs(doc, since) {
		a, err := e.getThreadmarkPost(u)
		if err != nil {
			return nil, err
		}

		if a != nil {
			articles = append(articles, *a)
		}
	}

	return articles, nil
}

func (e *Extractor) FindAll(f *kiasu.Feed) ([]kiasu.Post, error) {
	return e.FindSince(f, time.Time{})
}

func (e *Extractor) getThreadmarkURLs(doc *goquery.Document, since time.Time) []string {
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

func (e *Extractor) getThreadmarkPost(url string) (*kiasu.Post, error) {
	split := strings.Split(url, "#post-")

	var postID string
	if len(split) != 2 {
		postID = "first"
	} else {
		postID = split[1]
	}
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		return nil, err
	}

	var sel *goquery.Selection
	if postID == "first" {
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

	text, err := html2text.FromString(h)
	if err != nil {
		return nil, err
	}

	return &kiasu.Post{
		CreatedAt: postTime,
		Title:     title,
		Content:   text,
	}, nil
}
