package fanfictionnet

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/Puerkitobio/goquery"
	"github.com/fortytw2/hydrocarbon"
	"github.com/microcosm-cc/bluemonday"
)

var firstPost = "first"

// NewPlugin returns a fresh xenforo plugin
func NewPlugin() (*hydrocarbon.Plugin, error) {
	return &hydrocarbon.Plugin{
		Name:     "fanfiction",
		Configs:  configs,
		Validate: validate,
		Run:      run,
	}, nil
}

// list all configs up to the limit (scrapes for threads, basically)
func configs(c hydrocarbon.Client, p *hydrocarbon.Pagination) ([]hydrocarbon.Config, int, error) {
	return nil, 0, nil
}

// ensure a configuration is valid
func validate(c hydrocarbon.Client, cfg hydrocarbon.Config) error {
	return nil
}

// Run launches the given scrape and returns when it is finished
func run(c hydrocarbon.Client, cfg hydrocarbon.Config) ([]hydrocarbon.Post, error) {
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

	storyID, err := getStoryID(cfg.InitialURL)
	if err != nil {
		return nil, err
	}

	var articles []hydrocarbon.Post
	for i := 1; i < getNumberOfChapters(doc); i++ {
		a, err := getChapterPost(c, fmt.Sprintf("https://www.fanfiction.net/s/%s/%d", storyID, i))
		if err != nil {
			return nil, err
		}

		if a != nil {
			articles = append(articles, *a)
		}
	}

	return articles, nil
}

var reStoryID = regexp.MustCompile(`.*fanfiction.net/s/(\d+)/.*`)

func getStoryID(url string) (string, error) {
	submatch := reStoryID.FindAllStringSubmatch(url, 1)
	if len(submatch) != 1 {
		return "", errors.New("no story id found")
	}
	if len(submatch[0]) != 2 {
		return "", errors.New("no story id found")
	}
	return submatch[0][1], nil
}

func getNumberOfChapters(doc *goquery.Document) int {
	var num int
	doc.Find("#chap_select").First().Find("option").EachWithBreak(func(i int, sel *goquery.Selection) bool {
		num++
		return true
	})

	return num
}

func getChapterPost(c hydrocarbon.Client, url string) (*hydrocarbon.Post, error) {
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

	title := doc.Find("#chap_select").First().Text()
	content, err := doc.Find("#storytext").First().Html()
	if err != nil {
		return nil, err
	}

	p := bluemonday.UGCPolicy()
	html := p.Sanitize(content)

	return &hydrocarbon.Post{
		URL:      url,
		PostedAt: time.Time{},
		Title:    title,
		Content:  html,
	}, nil
}
