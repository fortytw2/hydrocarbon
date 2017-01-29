package pg

import (
	"fmt"
	"strings"

	"github.com/fortytw2/hydrocarbon"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// GetUser saves a user by ID
func (s *Store) GetUser(id string) (*hydrocarbon.User, error) {
	folderRows, err := s.db.Queryx(`
		WITH user_folders as (
			SELECT unnest(folder_ids) id
			FROM users WHERE
			id = $1
		) SELECT folders.name, folders.created_at, folders.updated_at, folders.id
		FROM folders
		INNER JOIN user_folders uf on folders.id=uf.id::uuid`, id)
	if err != nil {
		return nil, err
	}

	var folders []hydrocarbon.Folder
	for folderRows.Next() {
		var f hydrocarbon.Folder
		err = folderRows.StructScan(&f)
		if err != nil {
			return nil, err
		}

		var feedRows *sqlx.Rows
		feedRows, err = s.db.Queryx(`
			WITH folder_feeds as (
				SELECT unnest(feed_ids) id
				FROM folders WHERE
				id = $1
			) SELECT feeds.name, feeds.created_at, feeds.updated_at, feeds.id, feeds.plugin, feeds.initial_url
			FROM feeds
			INNER JOIN folder_feeds ff on feeds.id=ff.id::uuid`, f.ID)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		var feeds []hydrocarbon.Feed
		for feedRows.Next() {
			var fe hydrocarbon.Feed
			err = feedRows.StructScan(&fe)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			feeds = append(feeds, fe)
		}

		f.Feeds = feeds

		folders = append(folders, f)
	}

	userRow := s.db.QueryRowx("SELECT id, created_at, updated_at, analytics, email, encrypted_password, paid_until, active, confirmed, confirmation_token, token_created_at, stripe_customer_id FROM users WHERE id = $1", id)
	if userRow.Err() != nil {
		return nil, userRow.Err()
	}

	var u hydrocarbon.User
	err = userRow.StructScan(&u)
	if err != nil {
		return nil, err
	}
	u.Folders = folders

	return &u, nil
}

// GetUserByEmail gets a yser by email
func (s *Store) GetUserByEmail(email string) (*hydrocarbon.User, error) {
	row := s.db.QueryRowx("SELECT id, created_at, updated_at, stripe_customer_id, paid_until, analytics, email, encrypted_password, active, confirmed, confirmation_token, token_created_at FROM users WHERE email = $1", email)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var u hydrocarbon.User
	err := row.StructScan(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// SetStripeCustomerID sets the stripe ID for a given user
func (s *Store) SetStripeCustomerID(userID, stripeCustomerID string) error {
	_, err := s.db.Exec("UPDATE users SET stripe_customer_id = $1 WHERE id = $2;", stripeCustomerID, userID)
	return err
}

// AddFolder sets the stripe ID for a given user
func (s *Store) AddFolder(userID, folderID string) error {
	_, err := s.db.Exec("UPDATE users SET folder_ids = array_append(folder_ids, $1) WHERE id = $2;", folderID, userID)
	return err
}

// CreateUser saves a user and returns it, with it's new ID
func (s *Store) CreateUser(u *hydrocarbon.User) (*hydrocarbon.User, error) {
	row := s.db.QueryRowx(`
		INSERT INTO users (email, encrypted_password, analytics, active, confirmed, confirmation_token, token_created_at, folder_ids)
	    VALUES ($1, $2, $3, $4, $5, $6, $7, '{}')
		RETURNING id, created_at, updated_at, email, encrypted_password, active, confirmed, confirmation_token, token_created_at
	`, u.Email, u.EncryptedPassword, u.Analytics, u.Active, u.Confirmed, u.ConfirmationToken, u.TokenCreatedAt)
	if row.Err() != nil {
		if pqE := row.Err().(*pq.Error); strings.Contains(pqE.Message, "users_email_key") {
			return nil, hydrocarbon.ErrUserExists
		}

		return nil, row.Err()
	}

	var usr hydrocarbon.User
	err := row.StructScan(&usr)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}
