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
func (db *DB) CreateLoginToken(ctx context.Context, userID, userAgent, ip string) (string, error) {
	row := db.sql.QueryRowContext(ctx, `INSERT INTO login_tokens
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

// ActivateLoginToken activates the given LoginToken and returns the user
// the token was for
func (db *DB) ActivateLoginToken(ctx context.Context, token string) (string, error) {
	row := db.sql.QueryRowContext(ctx, `UPDATE login_tokens
										SET (used) = (true)
										WHERE token = $1
										AND expires_at > now()
										AND used = false
										RETURNING user_id;`, token)

	var userID string
	err := row.Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("no valid token found")
		}
		return "", err
	}

	return userID, nil
}

// CreateSession creates a new session for the user ID and returns the
// session key
func (db *DB) CreateSession(ctx context.Context, userID, userAgent, ip string) (string, error) {
	row := db.sql.QueryRowContext(ctx, `INSERT INTO sessions 
										(user_id, user_agent, ip)
										VALUES ($1, $2, $3)
										RETURNING key;`, userID, userAgent, ip)

	var key string
	err := row.Scan(&key)
	if err != nil {
		return "", err
	}

	return key, nil
}

// DeactivateSession invalidates the current session
func (db *DB) DeactivateSession(ctx context.Context, key string) error {
	_, err := db.sql.QueryContext(ctx, `UPDATE 
										sessions
										SET (active) = (false)
										WHERE key = $1;`, key)
	return err
}

// ActiveUserFromKey returns the active user ID from the session key
func (db *DB) ActiveUserFromKey(ctx context.Context, key string) (string, error) {
	row := db.sql.QueryRowContext(ctx, `SELECT user_id 
										FROM sessions 
										WHERE key = $1
										AND active = true;`, key)

	var userID string
	err := row.Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}
