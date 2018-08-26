package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/oklog/run"

	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/discollect"
	"github.com/fortytw2/hydrocarbon/gcs"
	"github.com/fortytw2/hydrocarbon/pg"
	"github.com/fortytw2/hydrocarbon/postmark"

	"github.com/fortytw2/hydrocarbon/plugins/fictionpress"
	"github.com/fortytw2/hydrocarbon/plugins/jsonfeed"
	"github.com/fortytw2/hydrocarbon/plugins/parahumans"
	"github.com/fortytw2/hydrocarbon/plugins/rss"

	"github.com/heroku/x/hmetrics"
)

func main() {
	var g run.Group

	var (
		autoExplain   = flag.Bool("autoexplain", false, "run EXPLAIN on every database query")
		noEmailVerify = flag.Bool("no-email-verify", false, "send login links in response to token request")
	)

	flag.Parse()

	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}

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
		imageDomain = os.Getenv("IMAGE_DOMAIN")
	} else {
		imageDomain = "http://localhost" + getPort("IMAGE_PORT", ":8082")
	}

	log.Println("hydrocarbon: launching api server on port", getPort("PORT", ":8080"), "for", domain)
	log.Println("hydrocarbon: launching image server on port", getPort("IMAGE_PORT", ":8082"), "for", imageDomain)

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

	var fs discollect.FileStore
	if gcpSA, ok := os.LookupEnv("GCP_SERVICE_ACCOUNT"); ok {
		gcpFS, err := gcs.NewFileStore(gcpSA, os.Getenv("IMAGE_BUCKET_NAME"))
		if err != nil {
			log.Fatal(err)
		}
		fs = gcpFS
	} else {
		localFS, err := discollect.NewLocalFS("./images", imageDomain+"/static/images/")
		if err != nil {
			log.Fatal(err)
		}
		fs = localFS

		imageH := &http.Server{
			Addr:    getPort("IMAGE_PORT", ":8082"),
			Handler: httpLogger(hydrocarbon.ErrorHandler(localFS.ServeHTTP), "hydrocarbon-images"),
		}
		{
			g.Add(imageH.ListenAndServe, func(error) {
				err := imageH.Shutdown(context.TODO())
				if err != nil && err != http.ErrServerClosed {
					log.Println("hydrocarbon: error shutting down http server", err)
				}
			})
		}
	}

	dc, err := discollect.New(
		// pg.DB is a discollect writer
		discollect.WithWriter(db),
		discollect.WithMetastore(db),
		discollect.WithFileStore(fs),
		discollect.WithPlugins(fictionpress.Plugin, parahumans.Plugin, rss.Plugin, jsonfeed.Plugin),
	)
	if err != nil {
		log.Fatal(err)
	}

	ua := hydrocarbon.NewUserAPI(db, ks, m, "hydrocarbon", stripePrivKey, paymentEnabled)
	if noEmailVerify != nil && *noEmailVerify {
		ua.DisableEmailVerification()
	}

	r := hydrocarbon.NewRouter(
		ua,
		hydrocarbon.NewFeedAPI(db, dc, ks),
		hydrocarbon.NewReadStatusAPI(db, ks),
		domain)

	h := &http.Server{
		Addr:    getPort("PORT", ":8080"),
		Handler: httpLogger(cspMiddleware(gziphandler.GzipHandler(r), imageDomain), "hydrocarbon-api"),
	}

	// if running on heroku, start reporting enhanced language metrics
	herokuMetrics()

	{
		g.Add(h.ListenAndServe, func(error) {
			err := h.Shutdown(context.TODO())
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
		log.Println(prefix+":", hydrocarbon.GetRemoteIP(req), req.Method, req.URL, elapsedTime)
	})
}

func cspMiddleware(router http.Handler, imageDomain string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Security-Policy", fmt.Sprintf(`default-src 'self' data:; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com data:; img-src 'self' data: %s`, imageDomain))
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		router.ServeHTTP(w, req)
	})
}

func herokuMetrics() {
	if os.Getenv("HEROKU_METRICS_URL") != "" {
		var logger hmetrics.ErrHandler = func(_ error) error { return nil }
		go func() {
			var backoff int64
			for backoff = 1; ; backoff++ {
				start := time.Now()

				hmetrics.Report(context.Background(), hmetrics.DefaultEndpoint, logger)
				if time.Since(start) > 5*time.Minute {
					backoff = 1
				}

				time.Sleep(time.Duration(backoff*10) * time.Second)
			}
		}()
	}
}
