package hydrocarbon

import (
	"context"
	"net/http"

	"github.com/bouk/httprouter"
)

// MachineStore is used to abstract machine->machine database calls
type MachineStore interface {
	UpdatePosts(ctx context.Context, feedID string, posts []*Post) error
}

// NewMachineRouter returns an *httprouter.Router suitable for PRIVATE
// communication only
func NewMachineRouter(ms MachineStore) *httprouter.Router {
	r := httprouter.New()

	r.PATCH("/posts", upsertPosts(ms))

	return r
}

func upsertPosts(ms MachineStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ?
		// ms.UpdatePosts(r.Context(), "idklol", nil)
	})
}
