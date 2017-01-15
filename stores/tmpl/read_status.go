package tmpl

import "github.com/fortytw2/hydrocarbon"

// GetReadStatus returns read status by ID
func (s *Store) GetReadStatus(id string) (*hydrocarbon.ReadStatus, error) {
	return nil, nil
}

// GetReadStatusByPostID returns read status for a given post
func (s *Store) GetReadStatusByPostID(postID, userID string) (*hydrocarbon.ReadStatus, error) {
	return nil, nil
}

// CreateReadStatus saves read status
func (s *Store) CreateReadStatus(rs *hydrocarbon.ReadStatus) (*hydrocarbon.ReadStatus, error) {
	return nil, nil
}
