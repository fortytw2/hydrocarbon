package hydrocarbon

import "time"

// PostStore provides primitives for storing and retrieving posts
type PostStore interface {
	GetPost(feedID, postID string) (*Post, error)
	SavePost(*Post) (*Post, error)
	GetPosts(feedID string, pg *Pagination) ([]Post, error)
}

// A Post is a single posting to a feed
type Post struct {
	ID     string `json:"id"`
	FeedID string `json:"feed_id"`

	CreatedAt   time.Time `json:"created_at"`
	RefreshedAt time.Time `json:"refreshed_at"`

	Title    string    `json:"title"`
	PostedAt time.Time `json:"posted_at"`
	URL      string    `json:"url"`

	Content string `json:"content"`
}
