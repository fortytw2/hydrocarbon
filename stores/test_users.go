package stores

import (
	"testing"

	"github.com/fortytw2/kiasu"
	"github.com/stretchr/testify/assert"
)

type userTest struct {
	Email    string
	Password string
}

func fuzz(len int) []userTest {
	var out []userTest
	for i := len; i > 0; i-- {
		out = append(out, userTest{
			Email:    randToken(8) + "@" + randToken(4) + ".com",
			Password: randToken(24),
		})
	}
	return out
}

// FuzzUserStore puts a lot of random users in the user store, then gets them back
func FuzzUserStore(t *testing.T, us kiasu.UserStore, n int) {
	for _, u := range fuzz(n) {
		newUser, err := us.SaveUser(&kiasu.User{
			Email:             u.Email,
			EncryptedPassword: u.Password,
		})
		if err != nil {
			if err == kiasu.ErrUserExists {
				continue
			}
			assert.Nil(t, err)
		}
		assert.NotEmpty(t, newUser)

		shouldMatch, err := us.GetUser(newUser.ID)
		assert.Nil(t, err)
		assert.Equal(t, shouldMatch.ID, newUser.ID)
		assert.Equal(t, u.Email, shouldMatch.Email)

		shouldMatchZwei, err := us.GetUserByEmail(u.Email)
		assert.Nil(t, err)
		assert.Equal(t, shouldMatchZwei.ID, newUser.ID)
		assert.Equal(t, u.Email, shouldMatchZwei.Email)
	}
}

// TestUserStore ensures a given userStore does what it should do
func TestUserStore(t *testing.T, us kiasu.UserStore) {
	u, err := us.SaveUser(&kiasu.User{
		Email:             "ian@ian.com",
		EncryptedPassword: "we12312312",
	})

	assert.Nil(t, err)
	assert.NotEmpty(t, u.ID)

	u2, err := us.GetUser(u.ID)

	assert.Nil(t, err)
	assert.Equal(t, u.ID, u2.ID)
}
