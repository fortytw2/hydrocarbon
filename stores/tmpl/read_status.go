package tmpl

import "github.com/fortytw2/kiasu"

// GetReadStatus returns read status by ID
func (s *Store) GetReadStatus(id string) (*kiasu.ReadStatus, error) {
	return nil, nil
}

// GetReadStatusByPostID returns read status for a given post
func (s *Store) GetReadStatusByPostID(postID, userID string) (*kiasu.ReadStatus, error) {
	return nil, nil
}

// SaveReadStatus saves read status
func (s *Store) SaveReadStatus(rs *kiasu.ReadStatus) (*kiasu.ReadStatus, error) {
	return nil, nil
}
