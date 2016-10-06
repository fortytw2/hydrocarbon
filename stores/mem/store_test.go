package mem

import (
	"testing"

	"github.com/fortytw2/kiasu/stores"
)

func TestMemStore(t *testing.T) {
	stores.Test(NewStore([]byte{1, 2, 3, 4, 4, 3}), t)
}
