package hydrocarbon

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fortytw2/dockertest"
)

func TestUserAPI(t *testing.T) {
	t.Parallel()

	container, err := dockertest.RunContainer("postgres:alpine", "5432", func(addr string) error {
		db, err := sql.Open("postgres", "postgres://postgres:postgres@"+addr+"?sslmode=disable")
		if err != nil {
			return err
		}

		return db.Ping()
	})
	defer container.Shutdown()
	if err != nil {
		t.Fatalf("could not start postgres, %s", err)
	}

	db, err := NewDB("postgres://postgres:postgres@" + container.Addr + "?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("create", userAPITestCreate(db))
}

func userAPITestCreate(db *DB) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		s := httptest.NewServer(http.HandlerFunc((&UserAPI{
			s: db,
			m: &MockMailer{},
		}).Register))

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
