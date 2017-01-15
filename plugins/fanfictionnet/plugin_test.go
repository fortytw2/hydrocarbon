package fanfictionnet

import (
	"net/http"
	"testing"
	"time"

	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/watney"
	"github.com/stretchr/testify/assert"
)

// goldenPostURLs at the time of recording the watney .har
var goldenPostURLs = []string{
	"https://www.fanfiction.net/s/5782108/1",
	"https://www.fanfiction.net/s/5782108/116",
}

func TestExtractor(t *testing.T) {
	tr := watney.Configure(http.DefaultTransport, t)
	c := &http.Client{
		Transport: tr,
	}
	defer watney.Save(c)

	p, err := NewPlugin()
	assert.Nil(t, err)

	// cfg, err := p.Configs(tC, 1)
	// assert.Nil(t, err)

	// err = p.Validate(tC, cfg[0])
	// assert.Nil(t, err)

	posts, err := p.Run(c, hydrocarbon.Config{
		InitialURL: "https://www.fanfiction.net/s/5782108/1/Harry-Potter-and-the-Methods-of-Rationality",
		Since:      time.Time{},
	})
	assert.Nil(t, err)

	assert.Equal(t, posts[0].URL, goldenPostURLs[0])
	assert.Equal(t, 121, len(posts))
}
