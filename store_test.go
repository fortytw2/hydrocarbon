package hydrocarbon_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	dockertest "gopkg.in/ory-am/dockertest.v2"

	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/internal/log"
	"github.com/fortytw2/hydrocarbon/stores/pg"
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

func TestCreateUser(t *testing.T) {
	var users = []struct {
		Valid    bool
		Dupe     bool
		Email    string
		Password string
	}{
		{true, false, "ian@fortytw2.com", "sa8dwu9djio23jl"},
		{false, false, "joe@barbados.com", "no"}, // invalid password
		{true, true, "ian@fortytw2.com", "sa8dwu9djio23jl"},
		{true, true, "ian@fortytw2.com", "sa8dwu9djio23jl"},
		{true, true, "ian@fortytw2.com", "sa8dwu9djio23jl"},
		{true, true, "ian@fortytw2.com", "sa8dwu9djio23jl"},
	}

	ps, err := pg.NewStore(log.NewNopLogger(), dsn)
	assert.Nil(t, err)

	s, err := hydrocarbon.NewStore(ps, []byte{1, 2, 3, 4, 2})
	assert.Nil(t, err)

	for _, u := range users {
		outU, err := s.CreateUser(u.Email, u.Password)
		if err != nil {
			if !u.Valid {
				assert.NotNil(t, err)
				continue
			} else if u.Dupe {
				assert.Equal(t, hydrocarbon.ErrUserExists, err)
				continue
			} else {
				t.Fatal(err)
			}
		}

		assert.Equal(t, outU.Email, u.Email)

		if u.Valid {
			assert.Nil(t, err)
		}
	}
}
