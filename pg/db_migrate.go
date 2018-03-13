package pg

import (
	"database/sql"
	"sort"
	"strings"
)

//go:generate go-bindata -pkg pg -mode 0644 -modtime 499137600 -o db_migrations_generated.go schema/

func runMigrations(db *sql.DB) error {
	err := verifyMigrationsTable(db)
	if err != nil {
		return err
	}

	count, err := countMigrations(db)
	if err != nil {
		return err
	}

	// assetNames returns the map keys, which are iterated in a random order
	// and not usable without sorting here if there are more than 1 migration.
	names := AssetNames()
	sort.Strings(names)

	for i, file := range names {
		// skip running ones we've clearly already ran
		if count > 0 {
			count--
			continue
		}

		migration := MustAsset(file)

		err := runMigration(i, migration, db)
		if err != nil {
			return err
		}

		cleanName := strings.TrimPrefix(file, "schema/")
		err = recordMigration(cleanName, db)
		if err != nil {
			return err
		}
	}

	return nil
}

func verifyMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS migrations (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		name TEXT NOT NULL UNIQUE
	);`)
	return err
}

func countMigrations(db *sql.DB) (int, error) {
	row := db.QueryRow(`SELECT count(*) FROM migrations;`)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func runMigration(num int, buf []byte, db *sql.DB) error {
	_, err := db.Exec(string(buf))
	return err
}

func recordMigration(name string, db *sql.DB) error {
	_, err := db.Query("INSERT INTO migrations (name) VALUES ($1);", name)
	return err
}
