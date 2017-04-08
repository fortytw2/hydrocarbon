package hydrocarbon

import "testing"

func TestMigrations(t *testing.T) {
	// spin up postgres

	db, err := NewDB("postgres://postgres:postgres@localhost:5432?sslmode=disable")
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
}
