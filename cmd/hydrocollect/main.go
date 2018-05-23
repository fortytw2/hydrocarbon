package main

import (
	"context"
	"log"
	"net/http"

	"github.com/fortytw2/hydrocarbon/discollect"
	"github.com/fortytw2/hydrocarbon/discollect/api"
	"github.com/oklog/run"

	"github.com/fortytw2/hydrocarbon/plugins/fictionpress"
	"github.com/fortytw2/hydrocarbon/plugins/parahumans"
)

func main() {
	dc, err := discollect.New(
		discollect.WithPlugins(fictionpress.Plugin, parahumans.Plugin),
	)
	if err != nil {
		log.Fatal(err)
	}

	// err = dc.LaunchScrape("fictionpress", &discollect.Config{
	// 	DynamicEntry: true,
	// 	Entrypoints:  []string{`https://www.fictionpress.com/s/2961893/1/Mother-of-Learning`},
	// 	Type:         "full",
	// 	Name:         "Mother of Learning",
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = dc.LaunchScrape("parahumans", &discollect.Config{
	// 	DynamicEntry: true,
	// 	Entrypoints:  []string{`https://parahumans.wordpress.com/2011/06/11/1-1`},
	// 	Type:         "full",
	// 	Name:         "Worm",
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	r := api.Router(dc)
	h := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	var g run.Group
	{
		g.Add(h.ListenAndServe, func(error) {
			h.Shutdown(context.Background())
		})
	}
	{
		g.Add(func() error { return dc.Start(1) }, func(error) {
			dc.Shutdown(context.Background())
		})
	}

	log.Fatal(g.Run())
}
