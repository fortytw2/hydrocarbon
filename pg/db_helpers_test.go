package pg

import (
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

func truncateTables(t *testing.T, db *DB) {
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

type testCase struct {
	name string
	do   func(t *testing.T) error
	want func(t *testing.T) error
}

func runCases(t *testing.T, db *DB, cases []testCase) {
	t.Helper()

	for _, tt := range cases {
		truncateTables(t, db)

		t.Run(tt.name, func(t *testing.T) {
			err := tt.do(t)
			if err != nil {
				t.Fatal(err)
			}

			if tt.want != nil {
				err = tt.want(t)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
