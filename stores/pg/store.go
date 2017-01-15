package pg

import (
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/internal/log"
	"github.com/fortytw2/hydrocarbon/internal/pgmigrate"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

//go:generate go-bindata -pkg pg -o migrations_generated.go schema/

// Store provides basic persistence primitives
type Store struct {
	db *sqlx.DB
}

// NewStore creates a primitive persistence layer
func NewStore(l log.Logger, dsn string) (hydrocarbon.PrimitiveStore, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	err = Migrate(l, db)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

// Migrate runs all migrations
func Migrate(l log.Logger, db *sqlx.DB) error {
	migrations, err := pgmigrate.LoadMigrations(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "schema"})
	if err != nil {
		return err
	}

	number, err := pgmigrate.DefaultConfig.Migrate(db.DB, migrations)
	l.Log("msg", "migrated postgres", "count", number.Len(), "total_migrations", migrations.Len())
	return err
}
