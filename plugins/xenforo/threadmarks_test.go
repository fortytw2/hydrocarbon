package xenforo

import (
	"testing"
	"time"

	"github.com/fortytw2/kiasu"
)

func TestExtractor(t *testing.T) {
	e := NewExtractor()

	_, err := e.FindSince(&kiasu.Feed{
		URL: "https://forums.spacebattles.com/threads/twinning-s-worm-altpower.408788/threadmarks",
	}, time.Time{})

	if err != nil {
		t.Fatal(err)
	}
}
