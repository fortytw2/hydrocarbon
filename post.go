package kiasu

import "time"

// A Post is a single posting to a feed
type Post struct {
	ID     int `json:"id"`
	FeedID int `json:"feed_id"`

	CreatedAt   time.Time `json:"created_at"`
	RefreshedAt time.Time `json:"refreshed_at"`

	Title    string    `json:"title"`
	PostedAt time.Time `json:"posted_at"`
	URL      string    `json:"url"`

	Content string `json:"content"`
}
