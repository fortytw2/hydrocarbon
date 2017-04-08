package hydrocarbon

import (
	"time"

	"github.com/lib/pq"
)

type FeedFolder struct {
}

type Feed struct {
	ID string `json:"id"`

	// DB Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// To track updates and refreshes of feeds
	LastRefreshedAt pq.NullTime `json:"last_refreshed_at"`
	LastEnqueuedAt  pq.NullTime `json:"last_enqueued_at"`

	// Plugin Name
	Plugin string `json:"plugin,omitempty"`
	// URL sent to the plugin
	URL string `json:"url"`
	// name for this feed
	Name        string `json:"name"`
	Description string `json:"description"`

	HexColor string `json:"hex_color"`
	IconURL  string `json:"icon_url"`

	UnreadCount int `json:"unread_count"`
}
