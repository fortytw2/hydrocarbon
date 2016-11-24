package stores

import (
	"testing"

	"github.com/fortytw2/kiasu"
	"github.com/stretchr/testify/assert"
)

// TestFeedStore ensures a given feedStore does what it should
func TestFeedStore(t *testing.T, fs kiasu.FeedStore) {
	f, err := fs.SaveFeed(&kiasu.Feed{
		Plugin:      "xenforo",
		Name:        "totally-test-forum",
		Description: "lol",
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, f.Name)
	assert.NotEmpty(t, f)

	badOutF, err := fs.GetFeed("potatosImNotaRealUUID")
	assert.NotNil(t, err)
	assert.Empty(t, badOutF)

	outF, err := fs.GetFeed(f.ID)
	assert.Nil(t, err)
	assert.NotEmpty(t, outF.Name)
	assert.Equal(t, outF.Name, f.Name)

	outFs, err := fs.GetFeeds(&kiasu.Pagination{Page: 0, PageSize: 10})
	assert.Nil(t, err)
	assert.Equal(t, len(outFs), 1)
}
