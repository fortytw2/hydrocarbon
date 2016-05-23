package kiasu

import "time"

// A Extractor is used to find articles from a given feed
type Extractor interface {
	// Validate determines if the extractor can parse the given feed
	Validate(*Feed) error

	Update(*Article) error
	FindSince(*Feed, time.Time) ([]Article, error)
}

