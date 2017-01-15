package pg

import "github.com/fortytw2/hydrocarbon"

// GetReadStatus returns read status by ID
func (s *Store) GetReadStatus(id string) (*hydrocarbon.ReadStatus, error) {
	row := s.db.QueryRowx("SELECT * FROM read_statuses WHERE id = $1", id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var rs hydrocarbon.ReadStatus
	err := row.StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

// GetReadStatusByPostID returns read status for a given post
func (s *Store) GetReadStatusByPostID(postID, userID string) (*hydrocarbon.ReadStatus, error) {
	row := s.db.QueryRowx("SELECT * FROM read_statuses WHERE post_id = $1 and user_id = $2", postID, userID)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var rs hydrocarbon.ReadStatus
	err := row.StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

// CreateReadStatus saves read status
func (s *Store) CreateReadStatus(rs *hydrocarbon.ReadStatus) (*hydrocarbon.ReadStatus, error) {
	row := s.db.QueryRowx(`
		INSERT INTO read_statuses (user_id, post_id, read_at, device_id, location)
	    VALUES ($1, $2, $3, $4, $5)
		RETURNING *
	`, rs.UserID, rs.PostID, rs.ReadAt, rs.Device, rs.Location)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var rss hydrocarbon.ReadStatus
	err := row.StructScan(&rss)
	if err != nil {
		return nil, err
	}

	return &rss, nil
}
