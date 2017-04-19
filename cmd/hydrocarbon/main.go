package main

import (
	"log"
	"net/http"
	"os"

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

	err = http.ListenAndServe(getPort(), hydrocarbon.NewRouter(hydrocarbon.NewUserAPI(db, &hydrocarbon.StdoutMailer{})))
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
