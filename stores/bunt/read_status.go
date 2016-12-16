package bunt

import (
	"encoding/json"

	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/internal/uuid"
	"github.com/tidwall/buntdb"
)

const (
	readStatusPrefix = "readstatus:"
)

// GetReadStatus returns read status by ID
func (s *Store) GetReadStatus(id string) (*hydrocarbon.ReadStatus, error) {
	var rs hydrocarbon.ReadStatus

	err := s.db.View(func(tx *buntdb.Tx) error {
		js, err := tx.Get(readStatusPrefix + id)
		if err != nil {
			return err
		}

		return json.Unmarshal([]byte(js), &rs)
	})
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

// SaveReadStatus saves read status
func (s *Store) SaveReadStatus(rs *hydrocarbon.ReadStatus) (*hydrocarbon.ReadStatus, error) {
	id := uuid.NewV4()
	rs.ID = id.String()

	buf, err := json.Marshal(rs)
	if err != nil {
		return nil, err
	}

	err = s.db.Update(func(tx *buntdb.Tx) error {
		_, _, err = tx.Set(readStatusPrefix+id.String(), string(buf), &buntdb.SetOptions{Expires: false})
		return err
	})
	if err != nil {
		return nil, err
	}

	return rs, nil
}

// GetReadStatusByPostID returns read status for a given post
func (s *Store) GetReadStatusByPostID(postID, userID string) (*hydrocarbon.ReadStatus, error) {
	return nil, nil
}
