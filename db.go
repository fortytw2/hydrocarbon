package hydrocarbon

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

var (
	ErrInvalidSession = errors.New("invalid session token")
)

// A DB is responsible for all interactions with postgres
type DB struct {
	sql *sql.DB
}

// NewDB returns a new database
func NewDB(dsn string) (*DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = runMigrations(db)
	if err != nil {
		return nil, err
	}

	return &DB{
		sql: db,
	}, nil
}

// CreateUser creates a new user and returns the users ID
func (db *DB) CreateUser(ctx context.Context, email string) (string, error) {
	row := db.sql.QueryRowContext(ctx, `INSERT INTO users 
							(email) 
							VALUES ($1)
							RETURNING id;`, email)

	var userID string
	err := row.Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

// CreateLoginToken creates a new one-time-use login token
func (db *DB) CreateLoginToken(ctx context.Context, userID string) (string, error) {
	row := db.sql.QueryRowContext(ctx, `INSERT INTO login_tokens 
							(user_id)
							VALUES ($1)
							RETURNING token;`, userID)

	var token string
	err := row.Scan(&token)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ActivateLoginToken activates the given LoginToken and returns the user
// the token was for
func (db *DB) ActivateLoginToken(ctx context.Context, token string) (string, error) {
	row := db.sql.QueryRowContext(ctx, `UPDATE login_tokens
							SET (used) = (true)
							WHERE token = $1
							AND expired_at > now()
							RETURNING user_id;`, token)

	var userID string
	err := row.Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

// CreateSession creates a new session for the user ID and returns the
// session token
func (db *DB) CreateSession(ctx context.Context, userID, userAgent, ip string) (string, error) {
	row := db.sql.QueryRowContext(ctx, `INSERT INTO sessions 
							(user_id, user_agent, ip)
							VALUES ($1, $2, $3)
							RETURNING token;`, userID, userAgent, ip)

	var token string
	err := row.Scan(&token)
	if err != nil {
		return "", err
	}

	return token, nil
}

// DeleteSession invalidates the current session
func (db *DB) DeleteSession(ctx context.Context, token string) error {
	_, err := db.sql.QueryContext(ctx, `UPDATE 
							sessions
							SET (active) = (false)
							WHERE token = $1;`, token)
	return err
}
