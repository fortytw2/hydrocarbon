package kiasu

import "time"

// A User is a registered (or not) user
type User struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	EncryptedPassword string `json:"encrypted_password"`
	OTPSecret         []byte `json:"otp_secret,omitempty"`

	LoginCount       int `json:"login_count"`
	FailedLoginCount int `json:"failed_login_count"`

	ActiveUntil          *time.Time `json:"active_until"`
	StripeSubscriptionID string     `json:"stripe_subscription_id"`

	Active            bool       `json:"active"`
	Confirmed         bool       `json:"confirmed"`
	ConfirmationToken *string    `json:"confirmation_token"`
	TokenCreatedAt    *time.Time `json:"token_created_at"`

	NotifyWindow   time.Duration `json:"notify_window"`
	LastNotifiedAt *time.Time    `json:"last_notified_at"`
}
