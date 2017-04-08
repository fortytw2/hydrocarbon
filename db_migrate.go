package hydrocarbon

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

//go:generate go-bindata -pkg hydrocarbon -mode 0644 -modtime 499137600 -o db_migrations_generated.go schema/

func runMigrations(db *sql.DB) error {
	err := verifyMigrationsTable(db)
	if err != nil {
		return err
	}

	count, err := countMigrations(db)
	if err != nil {
		return err
	}

	for i, file := range AssetNames() {
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

		cleanName := strings.TrimLeft("schema/", file)
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
	queries := strings.Split(string(buf), ";")

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for i, q := range queries {
		_, err := tx.Exec(q + ";")
		if err != nil {
			return fmt.Errorf("migrator: statement %d  in migration %d failed: %s", i, num, err)
		}
	}

	return tx.Commit()
}

func recordMigration(name string, db *sql.DB) error {
	_, err := db.Query("INSERT INTO migrations (name) VALUES ($1);", name)
	if err != nil {
		return err
	}

	return nil
}
