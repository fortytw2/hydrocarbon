package bunt

import (
	"github.com/fortytw2/kiasu"
	"github.com/tidwall/buntdb"
)

// Store provides basic persistence on top of buntdb
type Store struct {
	db *buntdb.DB
}

// NewMemStore creates a purely in-memory buntdb store
func NewMemStore() (kiasu.PrimitiveStore, error) {
	db, err := buntdb.Open(":memory:")
	if err != nil {
		return nil, err
	}

	err = setup(db)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

// NewStore creates an on-disk, persistent storage layer
func NewStore(filepath string) (kiasu.PrimitiveStore, error) {
	db, err := buntdb.Open(filepath)
	if err != nil {
		return nil, err
	}

	err = setup(db)
	if err != nil {
		return nil, err
	}

	err = db.SetConfig(buntdb.Config{
		SyncPolicy: buntdb.Always,
	})
	if err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

func setup(db *buntdb.DB) error {
	err := db.CreateIndex("user_email", "user:*", buntdb.IndexJSON("email"))
	if err != nil {
		if err == buntdb.ErrIndexExists {
			// all is good
		} else {
			return err
		}
	}

	return nil
}
