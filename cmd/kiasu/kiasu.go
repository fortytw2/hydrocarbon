package main

import (
	"net/http"
	"os"

	"github.com/fortytw2/kiasu/api"
	"github.com/fortytw2/kiasu/web"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/inconshreveable/log15.v2"
)

func main() {
	l := log15.New()
	l.Info("starting kiasu")

	port := os.Getenv("PORT")
	if port == "" {
		l.Crit("no port found")
		os.Exit(1)
	}

	r := httprouter.New()

	r.GET("/feed", web.ListFeeds(nil, nil))

	// r.GET("/v1/extractors", api.Extractors)
	// r.GET("/v1/feeds", api.Feeds)
	r.GET("/api/v1/feed/", api.ListFeeds(nil, nil))
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		l.Crit("could not start kiasu", "error", err.Error())
	}
}
