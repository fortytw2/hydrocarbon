package pg

import "github.com/fortytw2/hydrocarbon"

// GetPost returns a post by ID
func (s *Store) GetPost(postID string) (*hydrocarbon.Post, error) {
	row := s.db.QueryRowx("SELECT * FROM posts WHERE id = $1", postID)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var p hydrocarbon.Post
	err := row.StructScan(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// CreatePost saves a post
func (s *Store) CreatePost(post *hydrocarbon.Post) (*hydrocarbon.Post, error) {
	row := s.db.QueryRowx(`
		INSERT INTO posts (feed_id, posted_at, title, url, content)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *
	`, post.FeedID, post.PostedAt, post.Title, post.URL, post.Content)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var p hydrocarbon.Post
	err := row.StructScan(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// GetPosts paginates through all posts for a feed
func (s *Store) GetPosts(feedID string, pg *hydrocarbon.Pagination) ([]hydrocarbon.Post, error) {
	rows, err := s.db.Queryx("SELECT * FROM posts WHERE feed_id = $1 ORDER BY posted_at OFFSET $2 LIMIT $3", feedID, pg.Page, pg.PageSize)
	if err != nil {
		return nil, err
	}

	var posts []hydrocarbon.Post
	for rows.Next() {
		var tmpPost hydrocarbon.Post
		err := rows.StructScan(&tmpPost)
		if err != nil {
			return nil, err
		}
		posts = append(posts, tmpPost)
	}

	return posts, nil
}
