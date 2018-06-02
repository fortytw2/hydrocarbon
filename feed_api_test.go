//+build integration

package hydrocarbon_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/pg"
)

func feedApiTests(db *pg.DB) func(t *testing.T) {
	var cases = []apiCase{
		{
			name: "valid-create-feed",
			req: func(t *testing.T, baseDomain string) *http.Request {
				return httptest.NewRequest(http.MethodPost,
					baseDomain+"/v1/feed/create",
					bytes.NewBufferString(`{"name": "hc", "plugin": "ycombinators", "url": "https://ycombinator.com"}`))
			},
			resp: func(t *testing.T, mm *hydrocarbon.MockMailer, w *httptest.ResponseRecorder) {
				if w.Code != 200 {
					t.Fatal("did not return 200")
				}
				if !strings.Contains(w.Body.String(), `"id":`) {
					t.Fatal("did not appear to get an id back")
				}
			},
			authed: true,
		},
	}

	return func(t *testing.T) {
		runCases(t, db, cases)
	}
}
