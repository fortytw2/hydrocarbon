package main

import (
	"net/http"
	"os"

	"github.com/fortytw2/hydrocarbon/internal/log"
)

func main() {
	l := log.NewContext(log.NewLogfmtLogger(os.Stdout)).With("ts", log.DefaultTimestampUTC)
	l.Log("msg", "launching hydrocarbon", "port", getPort())

	// feeds := api.NewGroup("/feed")
	// feeds.GET("/", api.GetUserFeeds(l, db))
	// feeds.GET("/:id", api.GetSingleFeed(l, db))
	//
	// posts := api.NewGroup("/post")
	// posts.POST("/read", api.ReadPost(l, db))
	// posts.GET("/:id", api.GetSinglePost(l, db))

	err := http.ListenAndServe(getPort(), nil)
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
