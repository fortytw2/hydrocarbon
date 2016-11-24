package tmpl

import "github.com/fortytw2/kiasu"

// GetPost returns a post by ID
func (s *Store) GetPost(feedID, postID string) (*kiasu.Post, error) {
	return nil, nil
}

// SavePost saves a post
func (s *Store) SavePost(post *kiasu.Post) (*kiasu.Post, error) {
	return nil, nil
}

// GetPosts paginates through all posts for a feed
func (s *Store) GetPosts(feedID string, pg *kiasu.Pagination) ([]kiasu.Post, error) {
	return nil, nil
}
