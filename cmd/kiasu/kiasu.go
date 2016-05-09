package main

import (
	"net/http"

	"gopkg.in/inconshreveable/log15.v2"
	"github.com/julienschmidt/httprouter"
)

func main() {
	l := log15.New()
	l.Info("starting kiasu")

	r := httprouter.New()

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		l.Crit("could not start kiasu", "error", err.Error())
	}
}
