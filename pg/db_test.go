package pg

import (
	"context"
	"testing"

	"github.com/fortytw2/dockertest"
)

func setupTestDB(t *testing.T) (*DB, func()) {
	var db *DB

	container, err := dockertest.RunContainer("postgres:alpine", "5432", func(addr string) error {
		var err error
		db, err = NewDB("postgres://postgres:postgres@"+addr+"?sslmode=disable", false)
		return err
	})
	if err != nil {
		t.Fatalf("could not start postgres, %s", err)
	}

	return db, container.Shutdown
}

func TestUser(t *testing.T) {
	db, shutdown := setupTestDB(t)
	defer shutdown()

	t.Run("create", createUser(db))
	t.Run("defaultFolder", defaultFolder(db))
	t.Run("addFeed", addFeed(db))
}

func createUser(db *DB) func(t *testing.T) {
	return func(t *testing.T) {
		_, _, err := db.CreateOrGetUser(context.Background(), "ian@hydrocarbon.io")
		if err != nil {
			t.Fatalf("could not create user %s", err)
		}

		_, _, err = db.CreateOrGetUser(context.Background(), "ian@HYDroCARBon.io")
		if err != nil {
			t.Fatal("error on creating same user twice:", err)
		}
	}
}

func TestSession(t *testing.T) {
	db, shutdown := setupTestDB(t)
	defer shutdown()

	t.Run("create", createSession(db))
}

func createSession(db *DB) func(t *testing.T) {
	return func(t *testing.T) {
		id, _, err := db.CreateOrGetUser(context.Background(), "ian@createsession.io")
		if err != nil {
			t.Fatalf("could not create user %s", err)
		}

		_, _, err = db.CreateSession(context.Background(), id, "Firefox", "192.168.1.21")
		if err != nil {
			t.Fatalf("could not create session %s", err)
		}
	}
}

func defaultFolder(db *DB) func(t *testing.T) {
	return func(t *testing.T) {
		userID, _, err := db.CreateOrGetUser(context.TODO(), "ian@testpotatoes.rs")
		if err != nil {
			t.Fatal(err)
		}

		_, key, err := db.CreateSession(context.Background(), userID, "Firefox", "192.168.1.21")
		if err != nil {
			t.Fatalf("could not create session %s", err)
		}

		fid, err := db.getDefaultFolderID(context.Background(), key)
		if err != nil {
			t.Fatal(err)
		}

		if fid == "" {
			t.Fatal("no default folder id")
		}

		fid2, err := db.getDefaultFolderID(context.Background(), key)
		if err != nil {
			t.Fatal(err)
		}

		if fid2 != fid {
			t.Fatal("default folder creating many default folders")
		}
	}
}

func addFeed(db *DB) func(t *testing.T) {
	return func(t *testing.T) {
		userID, _, err := db.CreateOrGetUser(context.Background(), "fow2qe.awdwad@qdwad.com")
		if err != nil {
			t.Fatal("could not create user")
		}

		_, key, err := db.CreateSession(context.Background(), userID, "Firefox", "192.168.1.21")
		if err != nil {
			t.Fatalf("could not create session %s", err)
		}

		_, err = db.AddFeed(context.Background(), key, "", "testfeed", "testplugin", "https://www.goole.com")
		if err != nil {
			t.Fatal(err)
		}
	}
}
