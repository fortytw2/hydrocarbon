package hydrocarbon

import "time"

// context.Context keys for value propogation
var (
	UserContextKey = "user_context_key"
	AccessTokenKey = "access_token_context_key"
)

// UserStore handles storage and retrieval of users and their sessions
type UserStore interface {
	GetUser(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	SaveUser(*User) (*User, error)
}

// A User is a registered (or not) user
type User struct {
	CreatedAt time.Time `json:"created_at"`

	ID                string `json:"id"`
	Email             string `json:"email"`
	EncryptedPassword string `json:"encrypted_password"`

	LoginCount       int `json:"login_count"`
	FailedLoginCount int `json:"failed_login_count"`

	Active            bool      `json:"active"`
	Confirmed         bool      `json:"confirmed"`
	ConfirmationToken string    `json:"confirmation_token"`
	TokenCreatedAt    time.Time `json:"token_created_at"`

	NotifyWindow   time.Duration `json:"notify_window"`
	LastNotifiedAt time.Time     `json:"last_notified_at"`

	Feeds []Feed `json:"feeds"`
}
