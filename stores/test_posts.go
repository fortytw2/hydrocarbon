package stores

import (
	"testing"

	"github.com/fortytw2/hydrocarbon"
	"github.com/stretchr/testify/assert"
)

// TestPostStore ensures a given postStore does what it should
func TestPostStore(t *testing.T, ss hydrocarbon.PostStore) {
	_, err := ss.SavePost(&hydrocarbon.Post{
		FeedID:  "421",
		Title:   "Best Post Ever",
		Content: "1",
	})

	assert.Nil(t, err)
	_, err = ss.SavePost(&hydrocarbon.Post{
		FeedID:  "421",
		Title:   "Best Post Everest",
		Content: "2",
	})

	assert.Nil(t, err)
	lastPost, err := ss.SavePost(&hydrocarbon.Post{
		FeedID:  "421",
		Title:   "Best Post 4ever",
		Content: "3",
	})

	assert.Nil(t, err)

	p, err := ss.GetPost(lastPost.FeedID, lastPost.ID)
	assert.Nil(t, err)
	assert.Equal(t, p.Content, "3")

	posts, err := ss.GetPosts("421", &hydrocarbon.Pagination{Page: 0, PageSize: 10})
	assert.Nil(t, err)
	assert.Equal(t, 3, len(posts))
}
