package main

import (
	"log"
	"net/http"
	"os"

	"github.com/NYTimes/gziphandler"
	"github.com/fortytw2/hydrocarbon"
)

func main() {
	log.Println("starting hydrocarbon on port", getPort())
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("no postgres dsn found")
	}

	db, err := hydrocarbon.NewDB(dsn)
	if err != nil {
		log.Fatal("could not connect to postgres", err)
	}

	var redirect bool
	if redir := os.Getenv("HTTPS_REDIRECT"); redir != "" {
		redirect = true
	}

	r := hydrocarbon.NewRouter(hydrocarbon.NewUserAPI(db, &hydrocarbon.StdoutMailer{}), redirect)
	err = http.ListenAndServe(getPort(), gziphandler.GzipHandler(r))
	if err != nil {
		log.Fatal(err)
	}
}

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":8080"
}
