package bunt

import (
	"encoding/json"
	"time"

	"github.com/fortytw2/kiasu"
	"github.com/fortytw2/kiasu/internal/uuid"
	"github.com/tidwall/buntdb"
)

const (
	postPrefix = "post:"
)

// GetPost returns a post by ID
func (s *Store) GetPost(feedID, postID string) (*kiasu.Post, error) {
	var p kiasu.Post

	err := s.db.View(func(tx *buntdb.Tx) error {
		js, err := tx.Get(postPrefix + feedID + ":" + postID)
		if err != nil {
			return err
		}

		return json.Unmarshal([]byte(js), &p)
	})
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// SavePost saves a post
func (s *Store) SavePost(post *kiasu.Post) (*kiasu.Post, error) {
	err := s.db.CreateIndex("post_feed_id_"+post.FeedID, "post:"+post.FeedID+":*", buntdb.IndexJSON("feed_id"))
	if err != nil {
		if err == buntdb.ErrIndexExists {
			// all is good
		} else {
			return nil, err
		}
	}

	id := uuid.NewV4()
	post.ID = id.String()
	post.CreatedAt = time.Now()

	buf, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}

	err = s.db.Update(func(tx *buntdb.Tx) error {
		_, _, err = tx.Set(postPrefix+post.FeedID+":"+id.String(), string(buf), &buntdb.SetOptions{Expires: false})
		return err
	})
	if err != nil {
		return nil, err
	}

	return post, nil
}

// GetPosts paginates through all posts for a feed
func (s *Store) GetPosts(feedID string, pg *kiasu.Pagination) ([]kiasu.Post, error) {
	var limit = pg.PageSize

	var posts []kiasu.Post
	err := s.db.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("post_feed_id_"+feedID, func(key string, value string) bool {
			var p kiasu.Post
			err := json.Unmarshal([]byte(value), &p)
			if err != nil {
				return true
			}

			if limit != 0 {
				posts = append(posts, p)
				limit--
			}

			return true
		})

		return err
	})
	if err != nil {
		return nil, err
	}

	return posts, nil
}
