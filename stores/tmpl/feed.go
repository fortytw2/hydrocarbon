package tmpl

import "github.com/fortytw2/hydrocarbon"

// GetFeed returns a feed by its ID
func (s *Store) GetFeed(id string) (*hydrocarbon.Feed, error) {
	return nil, nil
}

// SaveFeed saves a feed and returns it with it's new ID
func (s *Store) SaveFeed(*hydrocarbon.Feed) (*hydrocarbon.Feed, error) {
	return nil, nil
}

// GetFeeds returns and filters on feeds
func (s *Store) GetFeeds(pg *hydrocarbon.Pagination) ([]hydrocarbon.Feed, error) {
	return nil, nil
}
