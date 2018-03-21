//+build integration

package hydrocarbon_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/pg"
)

func TestAPI(t *testing.T) {
	db, cancel := pg.SetupTestDB(t)
	defer cancel()

	t.Run("feed-api", feedApiTests(db))
}

type apiCase struct {
	name   string
	req    func(t *testing.T, basePath string) *http.Request
	resp   func(t *testing.T, em *hydrocarbon.MockMailer, w *httptest.ResponseRecorder)
	authed bool
}

func runCases(t *testing.T, db *pg.DB, cases []apiCase) {
	t.Helper()

	for _, tt := range cases {
		db.TruncateTables(t)

		mm := &hydrocarbon.MockMailer{}
		ks := hydrocarbon.NewKeySigner("test")
		h := hydrocarbon.NewRouter(
			hydrocarbon.NewUserAPI(db, ks, mm, "", "", false),
			hydrocarbon.NewFeedAPI(db, ks),
			"http://localhost:3000",
		)

		w := httptest.NewRecorder()

		t.Run(tt.name, func(t *testing.T) {
			req := tt.req(t, "http://localhost:3000")
			if tt.authed {
				ak := getAuthKey(t, db, ks)
				req.Header.Set("X-Hydrocarbon-Key", ak)
			}

			h.ServeHTTP(w, req)
			tt.resp(t, mm, w)
		})
	}
}

func getAuthKey(t *testing.T, db *pg.DB, ks *hydrocarbon.KeySigner) string {
	id, _, err := db.CreateOrGetUser(context.TODO(), "ian@hydrocarbon.io")
	if err != nil {
		t.Fatal(err)
	}

	_, key, err := db.CreateSession(context.TODO(), id, "test-ua", "192.168.1.254")
	if err != nil {
		t.Fatal(err)
	}

	signed, err := ks.Sign(key)
	if err != nil {
		t.Fatal(err)
	}
	return signed
}
