package pg

import (
	"strings"

	"github.com/fortytw2/hydrocarbon"
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
		LEFT JOIN user_folders uf on folders.id=uf.id::uuid`, id)
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
		folders = append(folders, f)
	}

	userRow := s.db.QueryRowx("SELECT id, created_at, updated_at, email, encrypted_password, failed_login_count, active, confirmed, confirmation_token, token_created_at FROM users WHERE id = $1", id)
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
	row := s.db.QueryRowx("SELECT * FROM users WHERE email = $1", email)
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

// CreateUser saves a user and returns it, with it's new ID
func (s *Store) CreateUser(u *hydrocarbon.User) (*hydrocarbon.User, error) {
	row := s.db.QueryRowx(`
		INSERT INTO users (email, encrypted_password, failed_login_count, active, confirmed, confirmation_token, token_created_at, folder_ids)
	    VALUES ($1, $2, $3, $4, $5, $6, $7, '{}')
		RETURNING id, created_at, updated_at, email, encrypted_password, failed_login_count, active, confirmed, confirmation_token, token_created_at
	`, u.Email, u.EncryptedPassword, u.FailedLoginCount, u.Active, u.Confirmed, u.ConfirmationToken, u.TokenCreatedAt)
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
