package pg

import (
	"testing"

	"github.com/fortytw2/dockertest"
)

const truncateAllTables string = `
CREATE OR REPLACE FUNCTION truncate_tables(_username text)
	RETURNS void AS
$func$
BEGIN
	RAISE NOTICE '%', 
	-- EXECUTE  -- dangerous, test before you execute!
	(SELECT 'TRUNCATE TABLE '
		|| string_agg(format('%I.%I', schemaname, tablename), ', ')
		|| ' CASCADE'
	FROM pg_tables
	WHERE tableowner = _username
	AND schemaname = 'public'
	);
END
$func$ LANGUAGE plpgsql;`

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

	_, err = db.sql.Exec(truncateAllTables)
	if err != nil {
		t.Fatalf("couldn't create truncate tables script")
	}

	return db, container.Shutdown
}

func truncateTables(t *testing.T, db *DB) {
	_, err := db.sql.Exec(`SELECT truncate_tables('postgres')`)
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
