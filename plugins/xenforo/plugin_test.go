package xenforo

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/fortytw2/kiasu"
	"github.com/fortytw2/watney"
	"github.com/stretchr/testify/assert"
)

// goldenPostURLs at the time of recording the watney .har
var goldenPostURLs = []string{
	"https://forums.spacebattles.com/threads/skein-worm-altpower-au.437953/",
	"https://forums.spacebattles.com/threads/skein-worm-altpower-au.437953/page-2#post-26434210",
	"https://forums.spacebattles.com/threads/skein-worm-altpower-au.437953/page-3#post-26667058",
	"https://forums.spacebattles.com/threads/skein-worm-altpower-au.437953/page-3#post-26761741",
	"https://forums.spacebattles.com/threads/skein-worm-altpower-au.437953/page-5#post-26948119",
	"https://forums.spacebattles.com/threads/skein-worm-altpower-au.437953/page-6#post-27048513",
	"https://forums.spacebattles.com/threads/skein-worm-altpower-au.437953/page-10#post-27453990",
	"https://forums.spacebattles.com/threads/skein-worm-altpower-au.437953/page-11#post-27676560",
	"https://forums.spacebattles.com/threads/skein-worm-altpower-au.437953/page-13#post-27988210",
}

func TestExtractor(t *testing.T) {
	tr := watney.Configure(http.DefaultTransport, t)
	c := &http.Client{
		Transport: tr,
	}
	defer watney.Save(c)

	p, err := NewPlugin()
	assert.Nil(t, err)

	// cfg, err := p.Configs(context.TODO(), tC, 1)
	// assert.Nil(t, err)

	// err = p.Validate(context.TODO(), tC, cfg[0])
	// assert.Nil(t, err)

	posts, err := p.Run(context.TODO(), c, kiasu.Config{
		InitialURL: "https://forums.spacebattles.com/threads/skein-worm-altpower-au.437953/threadmarks",
		Since:      time.Time{},
	})
	assert.Nil(t, err)

	for i, p := range posts {
		assert.Equal(t, p.URL, goldenPostURLs[i])
	}
}
