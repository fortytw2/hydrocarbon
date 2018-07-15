package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/oklog/run"

	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/discollect"
	"github.com/fortytw2/hydrocarbon/pg"
	"github.com/fortytw2/hydrocarbon/postmark"

	"github.com/fortytw2/hydrocarbon/plugins/fictionpress"
	"github.com/fortytw2/hydrocarbon/plugins/parahumans"
	"github.com/fortytw2/hydrocarbon/plugins/rss"
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

	db, err := pg.NewDB(dsn, *autoExplain)
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

	var imageDomain string
	if os.Getenv("IMAGE_DOMAIN") != "" {
		// assume port is OK
		domain = os.Getenv("IMAGE_DOMAIN")
	} else {
		domain = "http://localhost" + getPort("IMAGE_PORT", ":8082")
	}

	var cspDomains = domain + " " + imageDomain

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

	ks := hydrocarbon.NewKeySigner(signingKey)

	// enable stripe
	stripePrivKey, paymentEnabled := os.LookupEnv("STRIPE_PRIVATE_TOKEN")
	if paymentEnabled {
		log.Println("payment enabled, tokens required to login")
	} else {
		log.Println("payment not enabled, set STRIPE_PRIVATE_TOKEN to enable")
	}

	localFS, err := discollect.NewLocalFS("./images", imageDomain+"/static/images/")
	if err != nil {
		log.Fatal(err)
	}

	dc, err := discollect.New(
		// pg.DB is a discollect writer
		discollect.WithWriter(db),
		discollect.WithMetastore(db),
		discollect.WithFileStore(localFS),
		discollect.WithPlugins(fictionpress.Plugin, parahumans.Plugin, rss.Plugin),
	)
	if err != nil {
		log.Fatal(err)
	}

	r := hydrocarbon.NewRouter(hydrocarbon.NewUserAPI(db, ks, m, "hydrocarbon", stripePrivKey, paymentEnabled), hydrocarbon.NewFeedAPI(db, dc, ks), domain)

	h := &http.Server{
		Addr:    getPort("PORT", ":8080"),
		Handler: httpLogger(cspMiddleware(gziphandler.GzipHandler(r), cspDomains), "hydrocarbon-api"),
	}

	imageH := &http.Server{
		Addr:    getPort("IMAGE_PORT", ":8082"),
		Handler: httpLogger(hydrocarbon.ErrorHandler(localFS.ServeHTTP), "hydrocarbon-images"),
	}

	var g run.Group
	{
		g.Add(h.ListenAndServe, func(error) {
			err := h.Shutdown(context.TODO())
			if err != nil && err != http.ErrServerClosed {
				log.Println("hydrocarbon: error shutting down http server", err)
			}
		})
	}
	{
		g.Add(imageH.ListenAndServe, func(error) {
			err := imageH.Shutdown(context.TODO())
			if err != nil && err != http.ErrServerClosed {
				log.Println("hydrocarbon: error shutting down http server", err)
			}
		})
	}
	{
		g.Add(func() error {
			log.Println("launching scraper")
			return dc.Start(3)
		}, func(error) {
			log.Println("shutting down scraper")
			dc.Shutdown(context.Background())
		})
	}
	{
		g.Add(func() error {
			sigCh := make(chan os.Signal, 1)

			signal.Notify(sigCh, os.Interrupt)
			<-sigCh

			return errors.New("hydrocarbon: os initiated shutdown")
		}, func(error) {})
	}

	log.Fatal(g.Run())
}

func getPort(env string, def string) string {
	p := os.Getenv(env)
	if p != "" {
		return ":" + p
	}

	return def
}

func httpLogger(router http.Handler, prefix string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		router.ServeHTTP(w, req)
		finishTime := time.Now()
		elapsedTime := finishTime.Sub(startTime)
		log.Println(prefix+":", req.RemoteAddr, req.Method, req.URL, elapsedTime)
	})
}

func cspMiddleware(router http.Handler, hosts string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self' data: "+hosts)
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		router.ServeHTTP(w, req)
	})
}
