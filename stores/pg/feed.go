package pg

import "github.com/fortytw2/hydrocarbon"

// GetFeed returns a feed by its ID
func (s *Store) GetFeed(id string) (*hydrocarbon.Feed, error) {
	row := s.db.QueryRowx("SELECT * FROM feeds WHERE id = $1", id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var f hydrocarbon.Feed
	err := row.StructScan(&f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

// CreateFeed saves a feed and returns it with it's new ID
func (s *Store) CreateFeed(f *hydrocarbon.Feed) (*hydrocarbon.Feed, error) {
	row := s.db.QueryRowx(`
		INSERT INTO feeds (plugin, initial_url, name, description, hex_color, icon_url)
	    VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *
	`, f.Plugin, f.InitialURL, f.Name, f.Description, f.HexColor, f.IconURL)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var feed hydrocarbon.Feed
	err := row.StructScan(&feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}

// GetFeeds returns and filters on feeds
func (s *Store) GetFeeds(pg *hydrocarbon.Pagination) ([]hydrocarbon.Feed, error) {
	rows, err := s.db.Queryx("SELECT * FROM feeds OFFSET $1 LIMIT $2", pg.Page, pg.PageSize)
	if err != nil {
		return nil, err
	}

	var feeds []hydrocarbon.Feed
	for rows.Next() {
		var tmpFeed hydrocarbon.Feed
		err := rows.StructScan(&tmpFeed)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, tmpFeed)
	}

	return feeds, nil
}
