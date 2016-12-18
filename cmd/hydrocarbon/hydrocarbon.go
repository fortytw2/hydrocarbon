package main

import (
	"net/http"
	"os"

	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/internal/log"
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
