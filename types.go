package hydrocarbon

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// A Folder holds a collection of feeds
type Folder struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Feeds []*Feed `json:"feeds"`
}

// A Feed is a collection of posts
type Feed struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Title     string    `json:"title"`
	Plugin    string    `json:"plugin"`
	BaseURL   string    `json:"base_url"`

	Unread int `json:"unread"`

	Posts []*Post `json:"posts"`
}

// A Post is a single post on a feed
type Post struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	PostedAt  time.Time `json:"posted_at"`
	UpdatedAt time.Time `json:"updated_at"`

	OriginalURL string `json:"original_url"`

	Title  string `json:"title"`
	Author string `json:"author"`
	Body   string `json:"body"`

	Read bool `json:"read"`

	Extra map[string]interface{} `json:"extra"`
}

// ContentHash returns the stable hex encoded SHA256 of a post
func (p *Post) ContentHash() string {
	h := sha256.New()

	_, err := h.Write([]byte(fmt.Sprintf("%s:%s:%s", p.Title, p.Author, p.Body)))
	if err != nil {
		// pretty certain this cannot error
		panic(err)
	}

	return hex.EncodeToString(h.Sum(nil))
}

// A Session is a session
type Session struct {
	CreatedAt time.Time `json:"created_at"`
	UserAgent string    `json:"user_agent"`
	IP        string    `json:"ip"`
	Active    bool      `json:"active"`
}
