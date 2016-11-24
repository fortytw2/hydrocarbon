package bunt

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/fortytw2/kiasu/stores"
	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestMemStore(t *testing.T) {
	s, err := NewMemStore()
	assert.Nil(t, err)
	assert.NotNil(t, s)

	stores.TestAll(t, s)
}

func TestDiskStore(t *testing.T) {
	path := fmt.Sprintf("/tmp/bunt-%d", rand.Int63())

	s, err := NewStore(path)
	assert.Nil(t, err)
	assert.NotNil(t, s)

	stores.TestAll(t, s)

	err = os.RemoveAll(path)
	if err != nil {
		fmt.Println("could not clean up after test")
	}
}
