package hydrocarbon

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetRemoteIP(t *testing.T) {
	t.Parallel()

	{
		r, _ := http.NewRequest(http.MethodGet, "https://www.hydrocarbon.io/", nil)
		r.Header.Set("X-Real-IP", "193.167.12.23")

		ip := getRemoteIP(r)
		if ip != "193.167.12.23" {
			t.Error("x-real-ip is broken")
		}
	}

	{
		r, _ := http.NewRequest(http.MethodGet, "https://www.hydrocarbon.io/", nil)
		r.Header.Set("X-Forwarded-For", "193.167.12.23, 204.121.12.21")

		ip := getRemoteIP(r)
		if ip != "193.167.12.23" {
			t.Error("x-forwarded-for is broken")
		}
	}

	{
		r, _ := http.NewRequest(http.MethodGet, "https://www.hydrocarbon.io/", nil)
		r.RemoteAddr = "123.34.121.121"

		ip := getRemoteIP(r)
		if ip != "123.34.121.121" {
			t.Error("remote addr fallback is broken")
		}
	}
}

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
