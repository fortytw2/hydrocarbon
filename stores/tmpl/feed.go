package tmpl

import "github.com/fortytw2/kiasu"

// GetFeed returns a feed by its ID
func (s *Store) GetFeed(id string) (*kiasu.Feed, error) {
	return nil, nil
}

// SaveFeed saves a feed and returns it with it's new ID
func (s *Store) SaveFeed(*kiasu.Feed) (*kiasu.Feed, error) {
	return nil, nil
}

// GetFeeds returns and filters on feeds
func (s *Store) GetFeeds(pg *kiasu.Pagination) ([]kiasu.Feed, error) {
	return nil, nil
}
