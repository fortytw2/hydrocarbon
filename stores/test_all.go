package stores

import (
	"testing"

	"github.com/fortytw2/kiasu"
)

// TestAll is a helper function that tests all parts of a PrimitiveStore
func TestAll(t *testing.T, s kiasu.PrimitiveStore) {
	t.Run("user store", func(t *testing.T) {
		TestUserStore(t, s)
	})

	t.Run("read status store", func(t *testing.T) {
		TestReadStatusStore(t, s)
	})

	t.Run("feed store", func(t *testing.T) {
		TestFeedStore(t, s)
	})

	t.Run("post store", func(t *testing.T) {
		TestPostStore(t, s)
	})

	t.Run("session store", func(t *testing.T) {
		TestSessionStore(t, s)
	})

	t.Run("fuzz user store", func(t *testing.T) {
		FuzzUserStore(t, s, 256)
	})
}
