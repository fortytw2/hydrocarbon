package hydrocarbon

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserAPI(t *testing.T) {
	t.Parallel()

	db, shutdown := setupTestDB(t)
	defer shutdown()

	t.Run("create", userAPITestCreate(db))
}

func userAPITestCreate(db *DB) func(t *testing.T) {
	return func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc((&UserAPI{
			s: db,
			m: &MockMailer{},
		}).RequestToken))

		resp, err := http.Post(s.URL, "application/json", strings.NewReader(`{"email":"ian@hydrocarbon.io"}`))
		if err != nil {
			t.Fatalf("could not post test server %s", err)
		}

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(buf, registerSuccess) {
			t.Fatal("did not register account")
		}
	}
}
