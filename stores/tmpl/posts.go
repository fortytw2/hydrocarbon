package tmpl

import "github.com/fortytw2/hydrocarbon"

// GetPost returns a post by ID
func (s *Store) GetPost(feedID, postID string) (*hydrocarbon.Post, error) {
	return nil, nil
}

// SavePost saves a post
func (s *Store) SavePost(post *hydrocarbon.Post) (*hydrocarbon.Post, error) {
	return nil, nil
}

// GetPosts paginates through all posts for a feed
func (s *Store) GetPosts(feedID string, pg *hydrocarbon.Pagination) ([]hydrocarbon.Post, error) {
	return nil, nil
}
