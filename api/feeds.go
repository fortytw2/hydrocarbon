package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fortytw2/kiasu"
	"github.com/fortytw2/kiasu/extractor/xenforo"
	"github.com/fortytw2/kiasu/store"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/inconshreveable/log15.v2"
)

func ListFeeds(ds store.Store, l log15.Logger) httprouter.Handle {
	e := xenforo.NewExtractor()

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		articles, err := e.FindSince(&kiasu.Feed{
			BaseURL: r.URL.Query().Get("url"),
		}, time.Time{})
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(&articles)
	}
}
