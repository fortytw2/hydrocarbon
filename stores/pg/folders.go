package pg

import "github.com/fortytw2/hydrocarbon"

// CreateFolder saves a folder
func (s *Store) CreateFolder(folder *hydrocarbon.Folder) (*hydrocarbon.Folder, error) {
	row := s.db.QueryRowx(`
		INSERT INTO folders (name, feed_ids)
		VALUES ($1, '{}')
		RETURNING id, created_at, updated_at, name
	`, folder.Name)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var f hydrocarbon.Folder
	err := row.StructScan(&f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

// AddFeed adds a feed to a folder
func (s *Store) AddFeed(folderID, feedID string) error {
	_, err := s.db.Exec("UPDATE folders SET feed_ids = array_append(feed_ids, $1) WHERE id = $2;", feedID, folderID)
	return err
}
