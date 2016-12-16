package hydrocarbon_test

import (
	"testing"

	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/stores/bunt"
	"github.com/stretchr/testify/assert"
)

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

	ps, err := bunt.NewMemStore()
	assert.Nil(t, err)

	s, err := hydrocarbon.NewStore(ps, []byte{1, 2, 3, 4, 2})
	assert.Nil(t, err)

	for _, u := range users {
		outU, err := s.CreateUser(u.Email, u.Password)
		if err != nil {
			if !u.Valid {
				assert.NotNil(t, err)
				continue
			}
			if u.Dupe {
				assert.Equal(t, hydrocarbon.ErrUserExists, err)
				continue
			}
		}

		assert.Equal(t, outU.Email, u.Email)

		if u.Valid {
			assert.Nil(t, err)
		}
	}
}
