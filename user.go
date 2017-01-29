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
	CreateUser(*User) (*User, error)

	SetStripeCustomerID(userID, stripeID string) error
	AddFolder(userID, folderID string) error
}

// A User is a registered (or not) user
type User struct {
	ID string `json:"id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Analytics         bool   `json:"analytics"`
	Email             string `json:"email"`
	EncryptedPassword string `json:"encrypted_password"`

	Active            bool      `json:"active"`
	Confirmed         bool      `json:"confirmed"`
	ConfirmationToken string    `json:"confirmation_token"`
	TokenCreatedAt    time.Time `json:"token_created_at"`

	StripeCustomerID string `json:"stripe_customer_id"`
	PaidUntil        string `json:"paid_until"`

	Folders []Folder `json:"folders"`
}
