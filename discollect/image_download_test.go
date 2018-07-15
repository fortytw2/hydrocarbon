package discollect

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const onePxPng = `iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAAAAAA6fptVAAAACklEQVR4nGNiAAAABgADNjd8qAAAAABJRU5ErkJggg==`

func TestDownloadImages(t *testing.T) {
	imageServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		png, _ := base64.StdEncoding.DecodeString(onePxPng)
		w.Header().Set("Content-Type", "image/png")
		w.Write(png)
	}))
	defer imageServer.Close()

	stubby := NewStubFS()

	var cases = []struct {
		Name   string
		Input  string
		Output string
	}{
		{
			"simple",
			`<img src="` + imageServer.URL + `/img.png" />`,
			`<img src="` + stubby.URL + imageServer.URL + `/img.png"/>`,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			outText, err := DownloadImages(c.Input, http.DefaultClient, stubby)
			if err != nil {
				t.Error(err)
			}

			if !strings.Contains(outText, c.Output) {
				t.Errorf("output does not match expected, %s, %s", outText, c.Output)
			}
		})
	}

}
