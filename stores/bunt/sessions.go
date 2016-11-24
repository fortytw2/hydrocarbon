package bunt

import "github.com/fortytw2/kiasu"

// GetSession returns a session by ID
func (s *Store) GetSession(id string) (*kiasu.Session, error) {
	return nil, nil
}

// GetSessionsByUserID returns all sessions for a given user
func (s *Store) GetSessionsByUserID(userID string, pg *kiasu.Pagination) ([]kiasu.Session, error) {
	return nil, nil
}

// GetSessionByAccessToken returns the session by access token
func (s *Store) GetSessionByAccessToken(token string) (*kiasu.Session, error) {
	return nil, nil
}

// SaveSession saves a new session
func (s *Store) SaveSession(ses *kiasu.Session) (*kiasu.User, error) {
	return nil, nil
}
