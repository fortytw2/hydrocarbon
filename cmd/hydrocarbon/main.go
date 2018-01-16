package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/postmark"
)

func main() {
	var (
		autoExplain = flag.Bool("autoexplain", false, "run EXPLAIN on every database query")
	)

	flag.Parse()

	log.Println("starting hydrocarbon on port", getPort("PORT", ":8080"))
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("no postgres dsn found")
	}

	db, err := hydrocarbon.NewDB(dsn, *autoExplain)
	if err != nil {
		log.Fatal("could not connect to postgres", err)
	}

	var domain string
	if os.Getenv("DOMAIN") != "" {
		// assume port is OK
		domain = os.Getenv("DOMAIN")
	} else {
		domain = "http://localhost" + getPort("PORT", ":8080")
	}

	var m hydrocarbon.Mailer
	{
		if os.Getenv("POSTMARK_KEY") != "" {
			log.Println("sending mails via postmark")
			m = &postmark.Mailer{
				Key:    os.Getenv("POSTMARK_KEY"),
				Domain: domain,
				Client: http.DefaultClient,
			}

		} else {
			log.Println("sending mails to stdout")
			m = &hydrocarbon.StdoutMailer{Domain: domain}
		}
	}

	var signingKey string
	{
		if sk := os.Getenv("SIGNING_KEY"); sk != "" {
			log.Println("using signing key from env")
			signingKey = sk
		} else {
			log.Println("using default signing key, CHANGE ME IN PROD")
			signingKey = "DEV_SIGNING_KEY"
		}
	}

	// enable stripe
	stripePrivKey, paymentEnabled := os.LookupEnv("STRIPE_PRIVATE_TOKEN")
	if paymentEnabled {
		log.Println("payment enabled, tokens required to login")
	} else {
		log.Println("payment not enabled, set STRIPE_PRIVATE_TOKEN to enable")
	}

	ks := hydrocarbon.NewKeySigner(signingKey)

	r := hydrocarbon.NewRouter(hydrocarbon.NewUserAPI(db, ks, m, "hydrocarbon", stripePrivKey, paymentEnabled), hydrocarbon.NewFeedAPI(db, ks), domain)

	log.Println("serving private api on port", getPort("MACHINE_PORT", ":6060"))
	go http.ListenAndServe(getPort("MACHINE_PORT", ":6060"), httpLogger(gziphandler.GzipHandler(hydrocarbon.NewMachineRouter(db))))

	err = http.ListenAndServe(getPort("PORT", ":8080"), httpLogger(gziphandler.GzipHandler(r)))
	if err != nil {
		log.Fatal(err)
	}
}

func getPort(env string, def string) string {
	p := os.Getenv(env)
	if p != "" {
		return ":" + p
	}

	return def
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
