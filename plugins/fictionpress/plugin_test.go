//+build plugin_integration

package fictionpress

import (
	"net/http"
	"testing"

	"github.com/fortytw2/hydrocarbon/discollect"
)

func TestConfigValidator(t *testing.T) {
	feedTitle, _, err := Plugin.ConfigCreator("https://www.fanfiction.net/s/3401052/1/A-Black-Comedy", &discollect.HandlerOpts{
		Client: http.DefaultClient,
	})
	if err != nil {
		t.Fatal(err)
	}

	if feedTitle != "A Black Comedy" {
		t.Fatalf("got wrong feedTitle: '%s'", feedTitle)
	}
}
