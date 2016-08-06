package main

import (
	"net/http"
	"os"

	"github.com/fortytw2/kiasu/api"
	"github.com/go-kit/kit/log"
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	"golang.org/x/net/context"
)

func main() {
	l := log.NewContext(log.NewLogfmtLogger(os.Stdout)).With("ts", log.DefaultTimestampUTC)

	l.Log("msg", "launching kiasu", "port", os.Getenv("PORT"))

	r := xmux.New()

	v1 := r.NewGroup("/api/v1")
	// all routes for users
	api.AddUserRoutes(v1, nil, nil, nil)

	// feeds := api.NewGroup("/feed")
	// feeds.GET("/", api.GetUserFeeds(l, db))
	// feeds.GET("/:id", api.GetSingleFeed(l, db))
	//
	// posts := api.NewGroup("/post")
	// posts.POST("/read", api.ReadPost(l, db))
	// posts.GET("/:id", api.GetSinglePost(l, db))

	err := http.ListenAndServe(os.Getenv("PORT"), xhandler.New(context.Background(), r))
	if err != nil {
		l.Log("msg", "could not start kiasu", "error", err)
	}
}
