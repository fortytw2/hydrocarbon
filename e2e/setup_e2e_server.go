//+build integration

package e2e

import (
	"net/http"
	"testing"

	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/discollect"
)

func SetupE2EServer(t *testing.T, address string) (http.Server, *pg.DB, *hydrocarbon.MockMailer, func()) {
	t.Helper()

	db, cancel := pg.SetupTestDB()

	dc, _ := discollect.New(discollect.WithPlugins(&discollect.Plugin{
		Name:        "ycombinators",
		Entrypoints: []string{".*"},
		ConfigCreator: func(url string, ho *discollect.HandlerOpts) (string, *discollect.Config, error) {
			return "gotem", &discollect.Config{
				Type:        discollect.FullScrape,
				Entrypoints: []string{"gotem"},
			}, nil
		},
	}))

	mm := &hydrocarbon.MockMailer{}
	ks := hydrocarbon.NewKeySigner("test")
	h := hydrocarbon.NewRouter(
		hydrocarbon.NewUserAPI(db, ks, mm, "", "", false),
		hydrocarbon.NewFeedAPI(db, dc, ks),
		hydrocarbon.NewReadStatusAPI(db, ks),
		address,
	)

	return h, db, mm, cancel
}
