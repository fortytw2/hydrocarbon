package main

import (
	"net/http"

	"gopkg.in/inconshreveable/log15.v2"

	"github.com/fortytw2/kiasu/api"
	"github.com/julienschmidt/httprouter"
)

func main() {
	l := log15.New()
	l.Info("starting kiasu")

	r := httprouter.New()

	// r.GET("/v1/extractors", api.Extractors)
	// r.GET("/v1/feeds", api.Feeds)
	r.GET("/v1/feed/", api.ListFeeds(nil, nil))

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		l.Crit("could not start kiasu", "error", err.Error())
	}
}
