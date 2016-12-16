package bunt

import "github.com/fortytw2/hydrocarbon"

// GetSession returns a session by ID
func (s *Store) GetSession(id string) (*hydrocarbon.Session, error) {
	return nil, nil
}

// GetSessionsByUserID returns all sessions for a given user
func (s *Store) GetSessionsByUserID(userID string, pg *hydrocarbon.Pagination) ([]hydrocarbon.Session, error) {
	return nil, nil
}

// GetSessionByAccessToken returns the session by access token
func (s *Store) GetSessionByAccessToken(token string) (*hydrocarbon.Session, error) {
	return nil, nil
}

// SaveSession saves a new session
func (s *Store) SaveSession(ses *hydrocarbon.Session) (*hydrocarbon.User, error) {
	return nil, nil
}
