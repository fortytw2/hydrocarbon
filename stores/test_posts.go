package stores

import (
	"testing"

	"github.com/fortytw2/hydrocarbon"
	"github.com/stretchr/testify/assert"
)

// TestPostStore ensures a given postStore does what it should
func TestPostStore(t *testing.T, ps hydrocarbon.PrimitiveStore) {
	f, err := ps.CreateFeed(&hydrocarbon.Feed{
		Plugin:      "xenforo",
		Name:        "totally-test-forum-96",
		Description: "lol23",
		InitialURL:  "potatoes23",
		HexColor:    "xD",
		IconURL:     "http://potato.com/",
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, f.Name)
	assert.NotEmpty(t, f)

	_, err = ps.CreatePost(&hydrocarbon.Post{
		FeedID:  f.ID,
		Title:   "Best Post Ever",
		URL:     "https://potatoes.com",
		Content: "1",
	})

	assert.Nil(t, err)
	_, err = ps.CreatePost(&hydrocarbon.Post{
		FeedID:  f.ID,
		Title:   "Best Post Everest",
		URL:     "https://potatoes.co",
		Content: "2",
	})

	assert.Nil(t, err)
	lastPost, err := ps.CreatePost(&hydrocarbon.Post{
		FeedID:  f.ID,
		Title:   "Best Post 4ever",
		URL:     "https://potatoes.io",
		Content: "3",
	})

	assert.Nil(t, err)

	p, err := ps.GetPost(lastPost.ID)
	assert.Nil(t, err)
	assert.Equal(t, p.Content, "3")

	posts, err := ps.GetPosts(f.ID, &hydrocarbon.Pagination{Page: 0, PageSize: 10})
	assert.Nil(t, err)
	assert.Equal(t, 3, len(posts))
}
