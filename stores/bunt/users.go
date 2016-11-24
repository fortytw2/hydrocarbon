package bunt

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/fortytw2/kiasu"
	"github.com/fortytw2/kiasu/internal/uuid"
	"github.com/tidwall/buntdb"
)

const (
	userPrefix = "user:"
)

// GetUser saves a user by ID
func (s *Store) GetUser(id string) (*kiasu.User, error) {
	var u kiasu.User

	err := s.db.View(func(tx *buntdb.Tx) error {
		js, err := tx.Get(userPrefix + id)
		if err != nil {
			return err
		}

		return json.Unmarshal([]byte(js), &u)
	})
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// GetUserByEmail gets a user by their email
func (s *Store) GetUserByEmail(email string) (*kiasu.User, error) {
	var u kiasu.User
	err := s.db.View(func(tx *buntdb.Tx) error {
		// TODO: read buntdb docs
		err := tx.AscendGreaterOrEqual("user_email", email, func(key string, value string) bool {
			if !strings.Contains(value, email) {
				return true
			}

			var u2 kiasu.User
			err := json.Unmarshal([]byte(value), &u2)
			if err != nil {
				return true
			}

			if u2.Email == email {
				u = u2
				return false
			}

			return true
		})

		return err
	})
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// SaveUser saves a user and returns it, with it's new ID
func (s *Store) SaveUser(u *kiasu.User) (*kiasu.User, error) {
	id := uuid.NewV4()
	u.ID = id.String()
	u.CreatedAt = time.Now()

	buf, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}

	err = s.db.Update(func(tx *buntdb.Tx) error {
		_, _, err = tx.Set(userPrefix+id.String(), string(buf), &buntdb.SetOptions{Expires: false})
		return err
	})
	if err != nil {
		return nil, err
	}

	return u, nil
}
