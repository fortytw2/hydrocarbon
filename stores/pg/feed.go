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

// SetRefreshedAt updates refreshed at on a feed
func (s *Store) SetRefreshedAt(feedID string) error {
	_, err := s.db.Exec("UPDATE feeds SET last_refreshed_at = now() WHERE id = $1;", feedID)
	return err
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

// GetFeedsToUpdate returns a list of feeds that have either
// 1. never been refreshed
// or 2. not been refreshed in the last 5 minutes
func (s *Store) GetFeedsToUpdate(max int) ([]hydrocarbon.Feed, error) {
	rows, err := s.db.Queryx(`
		UPDATE feeds
		SET last_enqueued_at = now()
		FROM (
  			SELECT id
  			FROM   feeds
  			WHERE last_enqueued_at < (now() - interval '5 minutes')
				OR last_enqueued_at IS NULL
  			LIMIT $1
  			FOR UPDATE
  		) f
		WHERE feeds.id = f.id
		RETURNING *;`, max)
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
