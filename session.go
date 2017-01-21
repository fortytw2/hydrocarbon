package hydrocarbon

import (
	"time"

	"github.com/fortytw2/hydrocarbon/internal/token"
	"github.com/lib/pq"
)

// SessionStore stores sessions for users
type SessionStore interface {
	GetSession(id string) (*Session, error)
	GetSessionsByUserID(userID string, pg *Pagination) ([]Session, error)
	GetSessionByAccessToken(token string) (*Session, error)
	InvalidateSessionByToken(token string) error
	CreateSession(*Session) (*Session, error)
}

// A Session is a single user session
type Session struct {
	ID            string      `json:"id"`
	UserID        string      `json:"user_id"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	InvalidatedAt pq.NullTime `json:"invalidated_at"`
	ExpiresAt     time.Time   `json:"expires_at"`
	Token         string      `json:"token"`
}

// CreateToken creates a new random token and stores it in the session
func (s *Session) CreateToken() error {
	tok, err := token.GenerateRandomString(32)
	if err != nil {
		return err
	}

	s.Token = tok

	return nil
}
