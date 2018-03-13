package hydrocarbon

import (
	"context"
	"net/http"
)

// MachineStore is used to abstract machine->machine database calls
type MachineStore interface {
	UpdatePosts(ctx context.Context, feedID string, posts []*Post) error
}

// NewMachineRouter returns an *httprouter.Router suitable for PRIVATE
// communication only
func NewMachineRouter(ms MachineStore) http.Handler {
	fpr := &fixedPathRouter{
		paths: make(map[string]http.Handler),
	}

	fpr.paths["/posts"] = upsertPosts(ms)

	return fpr
}

func upsertPosts(ms MachineStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ?
		// ms.UpdatePosts(r.Context(), "idklol", nil)
	})
}
