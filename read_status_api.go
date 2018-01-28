package hydrocarbon

import "context"

// ReadStatusStore tracks read_statuses
type ReadStatusStore interface {
	MarkRead(ctx context.Context, postID, userID string)
}
