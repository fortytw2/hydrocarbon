//+build integration

package pg

import (
	"testing"

	"github.com/fortytw2/dockertest"
)

func SetupTestDB(t *testing.T) (*DB, func()) {
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

// TruncateTables is hidden behind build tags so we can use it in
// package hydrocarbon_test without fear of it being anywhere near a real build
func (db *DB) TruncateTables(t *testing.T) {
	// https://stackoverflow.com/a/12082038
	// TODO(fortytw2): this can be optimized by creating a template DB
	// and cloning / dropping it after running migrations
	_, err := db.sql.Exec(`
	DO
	$func$
	BEGIN
	EXECUTE (SELECT 'TRUNCATE TABLE ' || string_agg(oid::regclass::text, ', ') || ' CASCADE'
		FROM pg_class
		WHERE relkind = 'r'  -- only tables
		AND relnamespace = 'public'::regnamespace);
	END
	$func$;`)
	if err != nil {
		t.Fatal(err)
	}
}

type TestCase struct {
	name string
	do   func(t *testing.T) error
}

func RunCases(t *testing.T, db *DB, cases []TestCase) {
	t.Helper()

	for _, tt := range cases {
		db.TruncateTables(t)

		t.Run(tt.name, func(t *testing.T) {
			err := tt.do(t)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
