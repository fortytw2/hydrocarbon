package kiasu

import (
	"errors"
	"net/url"
)

// Errors defined for feed validation
var (
	ErrNoBaseURL = errors.New("no base url")
)

// A Feed is a single feed of articles - be it a spacebattles thread, a news
// site, or a FF story
type Feed struct {
	ID int `json:"-"`

	Name     string `json:"name"`
	BaseURL  string `json:"base_url"`
	IconURL  string `json:"icon_url"`
	HexColor string `json:"hex_color"`

	Extractor string `json:"extractor"`
}

// Validate performs two functions, canonicalizing the baseURL for deduplication
// of feeds, and ensuring the feed seems valid
func (f *Feed) Validate() error {
	if f.BaseURL == "" {
		return ErrNoBaseURL
	}

	u, err := url.Parse(f.BaseURL)
	if err != nil {
		return err
	}

	f.BaseURL = u.String()

	return nil
}
