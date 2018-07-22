package pg

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/fortytw2/hydrocarbon/discollect"

	"github.com/lib/pq"

	"github.com/fortytw2/hydrocarbon"
	"github.com/google/uuid"
	// postgres driver
	_ "github.com/lib/pq"
)

// A DB is responsible for all interactions with postgres
type DB struct {
	sql *sql.DB
}

// NewDB returns a new database
func NewDB(dsn string, autoExplain bool) (*DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = runMigrations(db)
	if err != nil {
		return nil, err
	}

	if autoExplain {
		_, err = db.Exec(`LOAD 'auto_explain';`)
		if err != nil {
			return nil, err
		}

		_, err = db.Exec(`SET auto_explain.log_min_duration = 0;`)
		if err != nil {
			return nil, err
		}
	}

	return &DB{
		sql: db,
	}, nil
}

// CreateOrGetUser creates a new user and returns the users ID
func (db *DB) CreateOrGetUser(ctx context.Context, email string) (string, bool, error) {
	row := db.sql.QueryRowContext(ctx, `
	INSERT INTO users 
	(email) 
	VALUES ($1)
	ON CONFLICT (email)
	DO UPDATE SET email = EXCLUDED.email
	RETURNING id, stripe_subscription_id;`, email)

	var userID string
	var stripeSubID sql.NullString
	err := row.Scan(&userID, &stripeSubID)
	if err != nil {
		return "", false, err
	}

	return userID, stripeSubID.Valid, nil
}

// SetStripeIDs sets a users stripe IDs
func (db *DB) SetStripeIDs(ctx context.Context, userID, customerID, subID string) error {
	_, err := db.sql.ExecContext(ctx, `
	UPDATE users 
	SET (stripe_customer_id, stripe_subscription_id) = ($1, $2)
	WHERE id = $3;`, customerID, subID, userID)

	return err
}

// CreateLoginToken creates a new one-time-use login token
func (db *DB) CreateLoginToken(ctx context.Context, userID, userAgent, ip string) (string, error) {
	row := db.sql.QueryRowContext(ctx, `
	INSERT INTO login_tokens
	(user_id, user_agent, ip)
	VALUES ($1, $2, $3::cidr)
	RETURNING token;`, userID, userAgent, ip)

	var token string
	err := row.Scan(&token)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ActivateLoginToken activates the given LoginToken and returns the user
// the token was for
func (db *DB) ActivateLoginToken(ctx context.Context, token string) (string, error) {
	row := db.sql.QueryRowContext(ctx, `
	UPDATE login_tokens
	SET used = true
	WHERE token = $1
	AND expires_at > now()
	AND used = false
	RETURNING user_id;`, token)

	var userID string
	err := row.Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("token invalid")
		}
		return "", err
	}

	return userID, nil
}

// CreateSession creates a new session for the user ID and returns the
// session key
func (db *DB) CreateSession(ctx context.Context, userID, userAgent, ip string) (email string, key string, err error) {
	row := db.sql.QueryRowContext(ctx, `
	INSERT INTO sessions 
	(user_id, user_agent, ip)
	VALUES ($1, $2, $3::cidr)
	RETURNING key;`, userID, userAgent, ip)
	err = row.Scan(&key)
	if err != nil {
		return "", "", err
	}

	row = db.sql.QueryRowContext(ctx, `
	SELECT email
	FROM users
	WHERE id = $1`, userID)
	err = row.Scan(&email)
	if err != nil {
		return "", "", err
	}

	return email, key, nil
}

// ListSessions lists all sessions a user has
func (db *DB) ListSessions(ctx context.Context, key string, page int) ([]*hydrocarbon.Session, error) {
	rows, err := db.sql.QueryContext(ctx, `
	SELECT created_at, user_agent, ip, active
	FROM sessions
	WHERE user_id = (SELECT user_id FROM sessions WHERE key = $1)
	LIMIT 25
	OFFSET $2`, key, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*hydrocarbon.Session
	for rows.Next() {
		var s hydrocarbon.Session
		err = rows.Scan(&s.CreatedAt, &s.UserAgent, &s.IP, &s.Active)
		if err != nil {
			return nil, err
		}
		out = append(out, &s)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return out, nil
}

// DeactivateSession invalidates the current session
func (db *DB) DeactivateSession(ctx context.Context, key string) error {
	_, err := db.sql.QueryContext(ctx, `
	UPDATE sessions
	SET (active) = (false)
	WHERE key = $1;`, key)

	return err
}

// AddFeed adds the given URL to the users default folder
// and links it across feed_folder
func (db *DB) AddFeed(ctx context.Context, sessionKey, folderID, title, plugin, feedURL string) (string, error) {
	if folderID == "" {
		// ensure we don't shadow folderID
		var err error
		folderID, err = db.getDefaultFolderID(ctx, sessionKey)
		if err != nil {
			return "", err
		}
	}

	tx, err := db.sql.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}

	row := tx.QueryRowContext(ctx, `
	INSERT INTO feeds
	(title, plugin, url)
	VALUES ($1, $2, $3)
	RETURNING id;`, title, plugin, feedURL)

	var feedID uuid.UUID
	err = row.Scan(&feedID)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return "", fmt.Errorf("%s - %s", err, txErr)
		}
		return "", err
	}

	_, err = tx.ExecContext(ctx, `
	INSERT INTO feed_folders
	(user_id, folder_id, feed_id)
	VALUES
	((SELECT user_id FROM sessions WHERE key = $1), $2, $3);`, sessionKey, folderID, feedID)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return "", fmt.Errorf("%s - %s", err, txErr)
		}
		return "", err
	}

	return feedID.String(), tx.Commit()
}

// getDefaultFolderID returns a users default folder ID
func (db *DB) getDefaultFolderID(ctx context.Context, sessionKey string) (string, error) {
	row := db.sql.QueryRowContext(ctx, `
	SELECT id FROM folders 
	WHERE name = 'default' 
	AND user_id = (SELECT user_id FROM sessions WHERE key = $1);`, sessionKey)

	var fid string
	err := row.Scan(&fid)
	if err != nil {
		// if there is no default folder, go create one
		if err == sql.ErrNoRows {
			row := db.sql.QueryRowContext(ctx, `
			INSERT INTO folders
			(user_id)
			VALUES 
			((SELECT user_id FROM sessions WHERE key = $1 LIMIT 1))
			RETURNING id;`, sessionKey)

			err = row.Scan(&fid)
			if err != nil {
				return "", fmt.Errorf("could not create default folder: %s", err)
			}

		} else {
			return "", fmt.Errorf("could not find default folder: %s", err)
		}

	}

	return fid, nil
}

// AddFolder creates a new folder
func (db *DB) AddFolder(ctx context.Context, sessionKey, name string) (string, error) {
	row := db.sql.QueryRow(`
	INSERT INTO folders 
	(user_id, name) 
	VALUES 
	((SELECT user_id FROM sessions WHERE key = $1), $2)
	RETURNING id;`, sessionKey, name)

	var id string
	err := row.Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

// RemoveFeed removes the given feed ID from the user
func (db *DB) RemoveFeed(ctx context.Context, sessionKey, folderID, feedID string) error {
	_, err := db.sql.ExecContext(ctx, `
	DELETE FROM feed_folders 
	WHERE user_id = (SELECT user_id FROM sessions WHERE key = $1 LIMIT 1)
	AND folder_id = $2
	AND feed_id = $3;`, sessionKey, folderID, feedID)

	return err
}

// GetFolders returns all of the folders for a user - if there are none it creates a
// default folder
func (db *DB) GetFolders(ctx context.Context, sessionKey string) ([]*hydrocarbon.Folder, error) {
	rows, err := db.sql.QueryContext(ctx, `
	SELECT fo.name as folder_name, fo.id as folder_id
	FROM folders fo
	WHERE fo.user_id = (SELECT user_id FROM sessions WHERE key = $1 LIMIT 1) 
	ORDER BY fo.created_at DESC;`, sessionKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	folders := make([]*hydrocarbon.Folder, 0)
	for rows.Next() {
		var folderName, folderID string

		err = rows.Scan(&folderName, &folderID)
		if err != nil {
			return nil, err
		}

		folders = append(folders, &hydrocarbon.Folder{
			ID:    folderID,
			Title: folderName,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return folders, nil
}

// GetFeedsForFolder returns a single feed
func (db *DB) GetFeedsForFolder(ctx context.Context, sessionKey, folderID string, limit, offset int) ([]*hydrocarbon.Feed, error) {
	rows, err := db.sql.QueryContext(ctx, `
	SELECT fe.id, fe.title, fe.url, fe.plugin
 	FROM feeds fe
	RIGHT JOIN feed_folders ff ON 
		(fe.id = ff.feed_id
		AND ff.user_id = (SELECT user_id FROM sessions WHERE key = $1) 
		AND ff.folder_id = $2)
	LIMIT $3 OFFSET $4;`, sessionKey, folderID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f []*hydrocarbon.Feed
	for rows.Next() {
		var feedID, feedTitle, feedURL, feedPlugin sql.NullString

		err = rows.Scan(&feedID, &feedTitle, &feedURL, &feedPlugin)
		if err != nil {
			return nil, err
		}

		if !feedID.Valid {
			continue
		}

		f = append(f, &hydrocarbon.Feed{
			ID:      feedID.String,
			Title:   feedTitle.String,
			Plugin:  feedPlugin.String,
			BaseURL: feedURL.String,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return f, nil
}

// GetFeed returns a single feed
func (db *DB) GetFeed(ctx context.Context, sessionKey, feedID string, limit, offset int) (*hydrocarbon.Feed, error) {
	rows, err := db.sql.QueryContext(ctx, `
	SELECT fe.id, fe.title, jsonb_agg(
		json_build_object('id', po.id, 'title', po.title, 'author', po.author, 'body', po.body, 'original_url', po.url, 'created_at', po.created_at, 'updated_at', po.updated_at, 'posted_at', po.posted_at)
	ORDER BY po.posted_at DESC) FILTER (WHERE po.id IS NOT NULL)
	FROM feeds fe
	LEFT JOIN posts po ON (fe.id = po.feed_id)
	WHERE fe.id = $2
	AND EXISTS (SELECT 1 FROM sessions WHERE key = $1)
	GROUP BY fe.id, fe.title
	LIMIT $3 OFFSET $4;`, sessionKey, feedID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feed := &hydrocarbon.Feed{
		ID:    feedID,
		Posts: make([]*hydrocarbon.Post, 0),
	}
	for rows.Next() {
		var feedID, feedTitle string
		var jsonBody []byte

		err := rows.Scan(&feedID, &feedTitle, &jsonBody)
		if err != nil {
			return nil, err
		}
		feed.Title = feedTitle

		if len(jsonBody) > 0 {
			err = json.Unmarshal(jsonBody, &feed.Posts)
			if err != nil {
				return nil, err
			}
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return feed, nil
}

// UpdatePosts upserts a smattering of posts into the db
func (db *DB) UpdatePosts(ctx context.Context, feedID string, posts []*hydrocarbon.Post) error {
	for _, p := range posts {
		var contentHash string
		err := db.sql.QueryRow(`
		INSERT INTO posts 
		(feed_id, content_hash, title, author, body, url, posted_at)
		VALUES 
		($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (content_hash) DO NOTHING
		RETURNING content_hash;`, feedID, p.ContentHash(), p.Title, p.Author, p.Body, p.OriginalURL, p.PostedAt).Scan(&contentHash)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			return err
		}
	}

	return nil
}

// Write saves off the post to the db
func (db *DB) Write(ctx context.Context, scrapeID uuid.UUID, f interface{}) error {
	hcp, ok := f.(*hydrocarbon.Post)
	if !ok {
		log.Println("did not get a post back from the scraper, skipping")
		return nil
	}

	_, err := db.sql.ExecContext(ctx, `
	INSERT INTO posts 
	(feed_id, content_hash, title, author, body, url, posted_at)
	VALUES 
	(
		(SELECT feed_id FROM scrapes WHERE id = $1), $2, $3, $4, $5, $6, $7
	)
	ON CONFLICT DO NOTHING;`, scrapeID, hcp.ContentHash(), hcp.Title, hcp.Author, hcp.Body, hcp.OriginalURL, hcp.PostedAt)
	return err
}

// Close implements io.Closer for pg.DB
func (db *DB) Close() error {
	return nil
}

// StartScrapes selects a subset of scrapes that should currently be running, but
// are not yet.
func (db *DB) StartScrapes(ctx context.Context, limit int) (ss []*StartedScrape, err error) {
	tx, err := db.sql.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	rollback := true
	// defer rollback if we throw an error
	defer func() {
		if rollback {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = fmt.Errorf("err: %s, rollbackErr: %s", err, rollbackErr)
			}
		}
	}()

	// FOR UPDATE SKIP LOCKED allows us to reduce contention against
	// any other instance running this same query at the same time.
	rows, err := tx.QueryContext(ctx, `
	SELECT id 
	FROM scrapes
	WHERE scheduled_start_at >= now()
	AND state = 'STOPPED'
	AND array_length(errors, 1) < 3
	LIMIT $1
	FOR UPDATE SKIP LOCKED;`, limit)
	if err != nil {
		return nil, err
	}

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	rows, err = tx.QueryContext(ctx, `
	UPDATE scrapes 
	SET state = 'RUNNING', started_at = now() 
	WHERE id = ANY($1)
	RETURNING id, feed_id, plugin, config;`, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	var ss []discollect.StartedScrape
	for rows.Next() {
		var s discollect.StartedScrape
		err = rows.Scan(&s.ID, &s.FeedID, &s.Plugin, &s.Config)
		if err != nil {
			return nil, err
		}
		ss = append(ss, s)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	rollback = false
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return ss, nil
}

// ListScrapes is used to list and filter scrapes, for both session resumption
// and UI purposes
func (db *DB) ListScrapes(ctx context.Context, statusFilter string) ([]*RunningScrape, error) {
	return nil, nil
}

// EndScrape marks a scrape as SUCCESS and records the number of datums and
// tasks returned
func (db *DB) EndScrape(ctx context.Context, id uuid.UUID, datums, tasks int) error {
	return nil
}

// ErrorScrape marks a scrape as ERRORED and adds the error to its list
func (db *DB) ErrorScrape(ctx context.Context, id uuid.UUID, err error) error {
	return nil
}
