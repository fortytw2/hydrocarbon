package hydrocarbon

import (
	"context"
	"testing"

	"github.com/fortytw2/dockertest"
)

func setupTestDB(t *testing.T) (*DB, func()) {
	var db *DB

	container, err := dockertest.RunContainer("postgres:alpine", "5432", func(addr string) error {
		var err error
		db, err = NewDB("postgres://postgres:postgres@" + addr + "?sslmode=disable")
		return err
	})
	if err != nil {
		t.Fatalf("could not start postgres, %s", err)
	}

	return db, container.Shutdown
}

func TestUser(t *testing.T) {
	t.Parallel()

	db, shutdown := setupTestDB(t)
	defer shutdown()

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
