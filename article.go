package kiasu

import "time"

// An Article is a single timestamped article. Can be a news peice, chapter of
// a book, or whatever
type Article struct {
	ID int `json:"id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ScrapedAt time.Time `json:"scraped_at"`

	Title   string `json:"title"`
	Author  string `json:"author"`
	Content string `json:"content"`
}
