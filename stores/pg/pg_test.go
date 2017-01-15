package pg

import (
	"fmt"
	"os"
	"testing"
	"time"

	dockertest "gopkg.in/ory-am/dockertest.v2"

	"github.com/fortytw2/hydrocarbon/internal/log"
	"github.com/fortytw2/hydrocarbon/stores"
	"github.com/stretchr/testify/assert"
)

var dsn string

func TestMain(m *testing.M) {
	if testing.Short() {
		fmt.Println("skipping pgmigrate test")
		return
	}

	c, err := dockertest.ConnectToPostgreSQL(1, 5*time.Second, func(u string) bool {
		dsn = u
		return true
	})
	if err != nil {
		fmt.Printf("Could not connect to database: %s", err)
		return
	}

	// Run tests
	result := m.Run()
	err = c.KillRemove()
	if err != nil {
		panic(err)
	}
	os.Exit(result)
}

func TestStore(t *testing.T) {
	s, err := NewStore(log.NewNopLogger(), dsn)
	assert.Nil(t, err)
	assert.NotNil(t, s)

	stores.TestUserStore(t, s)
	stores.TestReadStatusStore(t, s)
	stores.TestFeedStore(t, s)
	stores.TestPostStore(t, s)
	stores.TestSessionStore(t, s)
}
