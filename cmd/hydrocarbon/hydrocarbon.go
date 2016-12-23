package main

import (
	"net/http"
	"os"

	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/internal/log"
	"github.com/fortytw2/hydrocarbon/plugins/xenforo"
	"github.com/fortytw2/hydrocarbon/stores/bunt"
	"github.com/fortytw2/hydrocarbon/web"
)

func main() {
	l := log.NewContext(log.NewLogfmtLogger(os.Stdout)).With("ts", log.DefaultTimestampUTC)
	l.Log("msg", "launching hydrocarbon", "port", getPort())

	memStore, err := bunt.NewMemStore()
	if err != nil {
		l.Log("msg", "cannot start", "error", err)
		return
	}

	s, err := hydrocarbon.NewStore(memStore, []byte{1, 2, 3, 4})
	if err != nil {
		l.Log("msg", "cannot start", "error", err)
		return
	}

	go launchScraper(l, s)
	// go func() {
	// 	for {
	// 		feeds, err := s.Feeds.GetFeeds(&hydrocarbon.Pagination{
	// 			Page:     0,
	// 			PageSize: 10,
	// 		})
	// 		if err != nil {
	// 			panic(err)
	// 		}
	//
	// 		for _, f := range feeds {
	// 			l.Log("feed", f.ID, "name", f.Name)
	//
	// 			posts, _ := s.Posts.GetPosts(f.ID, &hydrocarbon.Pagination{
	// 				Page:     0,
	// 				PageSize: 10,
	// 			})
	//
	// 			for _, p := range posts {
	// 				l.Log("post", p.Title, "url", p.URL)
	// 			}
	// 		}
	//
	// 		time.Sleep(10 * time.Second)
	// 	}
	// }()

	r := web.Routes(s, l)
	err = http.ListenAndServeTLS(getPort(), "cert.pem", "key.pem", r)
	if err != nil {
		l.Log("msg", "could not start hydrocarbon", "error", err)
	}
}

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":8080"
}

func launchScraper(l log.Logger, s *hydrocarbon.Store) {
	_, err := s.Feeds.SaveFeed(&hydrocarbon.Feed{
		Plugin:      "xenforo",
		InitialURL:  "https://forums.spacebattles.com/threads/skein-worm-altpower-au.437953/threadmarks",
		Name:        "spacebattles-skein",
		Description: "lol",
	})
	if err != nil {
		panic(err)
	}

	plugins := map[string]hydrocarbon.Instantiator{
		"xenforo": xenforo.NewPlugin,
	}

	hydrocarbon.ScrapeLoop(l, s.Feeds, s.Posts, plugins)
}
