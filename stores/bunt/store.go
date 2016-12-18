package bunt

import (
	"github.com/fortytw2/hydrocarbon"
	"github.com/tidwall/buntdb"
)

// Store provides basic persistence on top of buntdb
type Store struct {
	db *buntdb.DB
}

// NewMemStore creates a purely in-memory buntdb store
func NewMemStore() (hydrocarbon.PrimitiveStore, error) {
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
func NewStore(filepath string) (hydrocarbon.PrimitiveStore, error) {
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
	err := db.Update(func(tx *buntdb.Tx) error {
		err := tx.CreateIndex("user_email", "user:*", buntdb.IndexJSON("email"))
		if err == buntdb.ErrIndexExists {
			// all is good
			return nil
		}
		return err
	})
	return err
}
