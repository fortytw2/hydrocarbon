package hydrocarbon

import (
	"context"
	"testing"

	"database/sql"

	"time"

	"github.com/fortytw2/hydrocarbon/testutil"
)

func TestMigrations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	addr, err := testutil.RunContainer(ctx, "postgres:alpine", "5432", func(addr string) error {
		db, err := sql.Open("postgres", "postgres://postgres:postgres@"+addr+"?sslmode=disable")
		if err != nil {
			return err
		}

		return db.Ping()
	})
	if err != nil {
		t.Fatalf("could not start postgres, %s", err)
	}

	// spin up postgres

	db, err := NewDB("postgres://postgres:postgres@" + addr + "?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	count, err := countMigrations(db.sql)
	if err != nil {
		t.Fatal(err)
	}

	// ensure we run all of the migrations
	if count != len(AssetNames()) {
		t.Fatalf("migrations not successful, expected %d but only found %d", len(AssetNames()), count)
	}

	cancel()
	time.Sleep(100 * time.Millisecond)
}
