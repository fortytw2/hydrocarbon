package hydrocarbon

import (
	"context"
	"database/sql"
	"testing"

	"github.com/fortytw2/dockertest"
)

func TestUser(t *testing.T) {
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

	t.Run("create", createUser(db))
}

func createUser(db *DB) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := db.CreateUser(context.Background(), "ian@hydrocarbon.io")
		if err != nil {
			t.Fatalf("could not create user %s", err)
		}

		_, err = db.CreateUser(context.Background(), "ian@hydrocarbon.io")
		if err == nil {
			t.Fatal("no error on creating same user twice")
		}
	}
}
