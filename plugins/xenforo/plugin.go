package xenforo

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Puerkitobio/goquery"
	"github.com/fortytw2/kiasu"
	"github.com/jaytaylor/html2text"
)

var firstPost = "first"

// NewPlugin returns a fresh xenforo plugin
func NewPlugin() (*kiasu.Plugin, error) {
	return &kiasu.Plugin{
		Name:     "xenforo",
		Configs:  configs,
		Validate: validate,
		Run:      run,
	}, nil
}

// list all configs up to the limit (scrapes for threads, basically)
func configs(c kiasu.Client, p *kiasu.Pagination) ([]kiasu.Config, int, error) {
	return nil, 0, nil
}

// ensure a configuration is valid
func validate(c kiasu.Client, cfg kiasu.Config) error {
	return nil
}

// Run launches the given scrape and returns when it is finished
func run(c kiasu.Client, cfg kiasu.Config) ([]kiasu.Post, error) {
	req, err := http.NewRequest("GET", cfg.InitialURL, nil)
	if err != nil {
		return nil, err
	}

	rsp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		return nil, err
	}
	err = rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	var articles []kiasu.Post
	for _, u := range getThreadmarkURLs(doc, cfg.Since) {
		a, err := getThreadmarkPost(c, u)
		if err != nil {
			return nil, err
		}

		if a != nil {
			articles = append(articles, *a)
		}
	}

	return articles, nil
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

func getThreadmarkPost(c kiasu.Client, url string) (*kiasu.Post, error) {
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

	rsp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rsp.Body.Close()
		if err != nil {
			fmt.Println("could not close rsp.Body", err)
		}
	}()

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
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

	text, err := html2text.FromString(h)
	if err != nil {
		return nil, err
	}

	return &kiasu.Post{
		URL:      url,
		PostedAt: postTime,
		Title:    title,
		Content:  text,
	}, nil
}
