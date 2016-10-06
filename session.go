package kiasu

import "time"

// A Session is a single user session
type Session struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Token     string    `json:"token"`
}
