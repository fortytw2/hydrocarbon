package kiasu

// A User is a single authenticated user
type User struct {
	ID int `json:"-"`

	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}
