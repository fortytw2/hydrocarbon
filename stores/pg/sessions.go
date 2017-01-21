package pg

import "github.com/fortytw2/hydrocarbon"

// GetSession returns a session by ID
func (s *Store) GetSession(id string) (*hydrocarbon.Session, error) {
	row := s.db.QueryRowx("SELECT * FROM sessions WHERE id = $1", id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var sess hydrocarbon.Session
	err := row.StructScan(&sess)
	if err != nil {
		return nil, err
	}

	return &sess, nil
}

// GetSessionsByUserID returns all sessions for a given user
func (s *Store) GetSessionsByUserID(userID string, pg *hydrocarbon.Pagination) ([]hydrocarbon.Session, error) {
	rows, err := s.db.Queryx("SELECT * FROM sessions WHERE user_id = $1 OFFSET $2 LIMIT $3", userID, pg.Page, pg.PageSize)
	if err != nil {
		return nil, err
	}

	var sess []hydrocarbon.Session
	for rows.Next() {
		var tempSess hydrocarbon.Session
		err := rows.StructScan(&tempSess)
		if err != nil {
			return nil, err
		}
		sess = append(sess, tempSess)
	}

	return sess, nil
}

// GetSessionByAccessToken returns the session by access token
func (s *Store) GetSessionByAccessToken(token string) (*hydrocarbon.Session, error) {
	row := s.db.QueryRowx("SELECT * FROM sessions WHERE token = $1 AND invalidated_at IS NULL AND expires_at > now();", token)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var sess hydrocarbon.Session
	err := row.StructScan(&sess)
	if err != nil {
		return nil, err
	}

	return &sess, nil
}

// InvalidateSessionByToken invalidates a given session
func (s *Store) InvalidateSessionByToken(token string) error {
	_, err := s.db.Exec("UPDATE sessions SET invalidated_at = now() WHERE token = $1;", token)
	return err
}

// CreateSession saves a new session
func (s *Store) CreateSession(ses *hydrocarbon.Session) (*hydrocarbon.Session, error) {
	row := s.db.QueryRowx(`
		INSERT INTO sessions (user_id, expires_at, token)
	    VALUES ($1, $2, $3)
		RETURNING *
	`, ses.UserID, ses.ExpiresAt, ses.Token)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var sess hydrocarbon.Session
	err := row.StructScan(&sess)
	if err != nil {
		return nil, err
	}

	return &sess, nil
}
