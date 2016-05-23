package spacebattles

import (
	"testing"
	"time"

	"github.com/fortytw2/kiasu"
)

func TestExtractor(t *testing.T) {
	e := NewExtractor()

	_, err := e.FindSince(&kiasu.Feed{
		BaseURL: "https://forums.spacebattles.com/threads/a-skittering-heart-worm-kingdom-hearts.371816/threadmarks",
	}, time.Now().AddDate(0, -12, 0))

	if err != nil {
		t.Fatal(err)
	}
}
