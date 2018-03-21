//+build integration

package pg

import (
	"testing"
)

func TestMigrations(t *testing.T) {
	db, shutdown := SetupTestDB(t)
	defer shutdown()

	count, err := countMigrations(db.sql)
	if err != nil {
		t.Fatal(err)
	}

	// ensure we run all of the migrations
	if count != len(AssetNames()) {
		t.Fatalf("migrations not successful, expected %d but only found %d", len(AssetNames()), count)
	}
}
