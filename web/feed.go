package web

import (
	"net/http"
	"text/template"
	"time"

	"github.com/fortytw2/kiasu"
	"github.com/fortytw2/kiasu/extractor/spacebattles"
	"github.com/fortytw2/kiasu/store"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/inconshreveable/log15.v2"
)

func ListFeeds(ds store.Store, l log15.Logger) httprouter.Handle {
	e := spacebattles.NewExtractor()

	t, _ := template.New("feed").Parse(`
	<html>
		<head>
			<title>Kiasu</title>
			<style>

p {
    white-space: pre-wrap;
}
			</style>

		</head>

		<body>
			{{range .}} 
			<h3>{{.Title}}</h3>
			<p> {{.Content}}</p> 

			{{end}}
		</body>
	</html>

		`)

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		articles, err := e.FindSince(&kiasu.Feed{
			BaseURL: r.URL.Query().Get("url"),
		}, time.Time{})
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t.Execute(w, articles)

	}
}
