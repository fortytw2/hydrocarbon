package hydrocarbon

import (
	"context"
	"net/http"
)

// ReadStatusStore tracks read_statuses
type ReadStatusStore interface {
	MarkRead(ctx context.Context, postID, sessionKey string) error
}

type ReadStatusAPI struct {
	s  ReadStatusStore
	ks *KeySigner
}

// NewReadStatusAPI returns a new Feed API
func NewReadStatusAPI(s ReadStatusStore, ks *KeySigner) *ReadStatusAPI {
	return &ReadStatusAPI{
		s:  s,
		ks: ks,
	}
}

// MarkRead marks the given post as read
func (rs *ReadStatusAPI) MarkRead(w http.ResponseWriter, r *http.Request) error {
	key, err := rs.ks.Verify(r.Header.Get("X-Hydrocarbon-Key"))
	if err != nil {
		return err
	}

	var readReq struct {
		PostID string `json:"post_id"`
	}

	err = limitDecoder(r, &readReq)
	if err != nil {
		return err
	}

	err = rs.s.MarkRead(r.Context(), key, readReq.PostID)
	if err != nil {
		return err
	}

	return writeSuccess(w, map[string]bool{
		readReq.PostID: true,
	})
}
