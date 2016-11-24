package tmpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	s, err := NewStore()
	assert.Nil(t, err)
	assert.NotNil(t, s)

	// stores.TestUserStore(t, s)
	// stores.TestReadStatusStore(t, s)
	// stores.TestFeedStore(t, s)
	// stores.TestPostStore(t, s)
	// stores.TestSessionStore(t, s)
}
