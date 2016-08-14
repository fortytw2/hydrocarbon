package main

import (
	"net/http"
	"os"

	"github.com/fortytw2/kiasu/api"
	"github.com/go-kit/kit/log"
)

func main() {
	l := log.NewContext(log.NewLogfmtLogger(os.Stdout)).With("ts", log.DefaultTimestampUTC)
	l.Log("msg", "launching kiasu", "port", os.Getenv("PORT"))

	// feeds := api.NewGroup("/feed")
	// feeds.GET("/", api.GetUserFeeds(l, db))
	// feeds.GET("/:id", api.GetSingleFeed(l, db))
	//
	// posts := api.NewGroup("/post")
	// posts.POST("/read", api.ReadPost(l, db))
	// posts.GET("/:id", api.GetSinglePost(l, db))

	err := http.ListenAndServe(":"+os.Getenv("PORT"), api.Routes(l, nil, nil))
	if err != nil {
		l.Log("msg", "could not start kiasu", "error", err)
	}
}
