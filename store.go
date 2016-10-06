package kiasu

import (
	"time"

	"context"
)

// Store is responsible for persistent (or not) data storage and retrieval
type Store interface {
	UserStore
	FeedStore
	PluginStore
	PaymentStore
}

// UserStore handles storage and retrieval of users and their sessions
type UserStore interface {
	// GetUser returns a users details
	GetUser(ctx context.Context, accessToken string) (*User, error)
	// CreateUser creates a new user, session and returns a confirmation_token
	// given email + encrypted password
	CreateUser(ctx context.Context, m Mailer, email string, pw string) (string, error)
	ActivateUser(ctx context.Context, confirmToken string) (string, error)
	// NewSession creates a new session
	NewSession(ctx context.Context, email string, pw string) (string, error)

	// Session management
	GetActiveSessions(ctx context.Context, accessToken string, p *Pagination) ([]Session, error)
	GetPastSessions(ctx context.Context, accessToken string, p *Pagination) ([]Session, error)
	InvalidateToken(ctx context.Context, accessToken string) error
}

// PaymentStore handles payment information
type PaymentStore interface {
	Charge(ctx context.Context, accessToken, chargeToken string) error
	// GetUsersByExpiry returns users sorted by and filtered by expiry date
	GetUsersByExpiry(ctx context.Context, m Mailer, expireAfter time.Time, p *Pagination) ([]User, error)
	// AddSubscription is used to add a subscription to a user with the given email
	AddSubscription(ctx context.Context, email string, activeUntil time.Time) error
}

// FeedStore handles storage and retrieval of feeds
type FeedStore interface {
	GetFeeds(ctx context.Context, p *Pagination) ([]Feed, error)
	// GetUserFeeds returns a users feeds, ordered by their SortOrder
	GetUserFeeds(ctx context.Context, accessToken string, p *Pagination) ([]Feed, error)
	// ReOrderFeed allows the ordering of a feed to be changed, moving all
	// feeds _down_ one from the change
	ReOrderFeed(ctx context.Context, accessToken string, feedID string, newOrder int) ([]Feed, error)

	GetFeedPosts(ctx context.Context, accessToken string, feedID string, p *Pagination) ([]Post, error)
}

// PluginStore handles storage and retrieval of both inproc and rpc plugins
type PluginStore interface {
	GetPlugins(ctx context.Context, p *Pagination, active bool) ([]Plugin, error)
	GetUserPlugins(ctx context.Context, accessToken string) ([]Plugin, error)
	SearchPlugins(ctx context.Context, accessToken, query string, p *Pagination) ([]Plugin, error)

	GetPluginStatus(ctx context.Context, accessToken, pluginID int) ([]Healthcheck, error)

	RegisterInProcPlugin(ctx context.Context, pl Plugin, title, desc string) error
	RegisterRPCPlugin(ctx context.Context, accessToken, url, title, desc string) error
}
