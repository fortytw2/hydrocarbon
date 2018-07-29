//+build plugin_integration

package fictionpress

import (
	"net/http"
	"testing"

	"github.com/fortytw2/hydrocarbon/discollect"
)

func TestConfigValidator(t *testing.T) {
	exConfig := &discollect.Config{
		Name:         "full",
		Entrypoints:  []string{"https://www.fanfiction.net/s/3401052/1/A-Black-Comedy"},
		DynamicEntry: true,
	}

	feedTitle, err := Plugin.ConfigValidator(&discollect.HandlerOpts{
		Client: http.DefaultClient,
		Config: exConfig,
	})
	if err != nil {
		t.Fatal(err)
	}

	if feedTitle != "A Black Comedy" {
		t.Fatalf("got wrong feedTitle- '%s'", feedTitle)
	}
}
