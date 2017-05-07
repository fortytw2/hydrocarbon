package main

import (
	"log"
	"net/http"
	"os"
	"time"

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

	var domain string
	if os.Getenv("DOMAIN") != "" {
		// assume port is OK
		domain = os.Getenv("DOMAIN")
	} else {
		domain = "http://localhost" + getPort()
	}
	log.Println("ui will target", domain+"/api", "for api requests")

	sentryPublic := os.Getenv("SENTRY_PUBLIC_DSN")
	log.Println("using SENTRY_PUBLIC_DSN", sentryPublic)

	r := hydrocarbon.NewRouter(hydrocarbon.NewUserAPI(db, &hydrocarbon.StdoutMailer{Domain: domain}), domain, sentryPublic)
	err = http.ListenAndServe(getPort(), httpLogger(gziphandler.GzipHandler(r)))
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

func httpLogger(router http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		router.ServeHTTP(w, req)
		finishTime := time.Now()
		elapsedTime := finishTime.Sub(startTime)
		log.Println("hydrocarbon:", req.Method, req.URL, elapsedTime)
	})
}
