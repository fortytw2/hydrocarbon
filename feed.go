package kiasu

import "time"

// FeedStore saves and persists feeds and posts
type FeedStore interface {
	GetFeed(id string) (*Feed, error)
	SaveFeed(*Feed) (*Feed, error)
	GetFeeds(pg *Pagination) ([]Feed, error)
}

// A Feed is a single encapsulating unit around a source of news / posts / whatever
type Feed struct {
	ID     string `json:"id"`
	Plugin string `json:"plugin,omitempty"`

	CreatedAt   time.Time  `json:"created_at"`
	RefreshedAt *time.Time `json:"refreshed_at"`

	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`

	HexColor string `json:"hex_color"`
	IconURL  string `json:"icon_url"`
}
