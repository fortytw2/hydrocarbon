package kiasu

import "time"

// ReadStatusStore stores and retrieves user read statuses
type ReadStatusStore interface {
	GetReadStatus(id string) (*ReadStatus, error)
	GetReadStatusByPostID(postID, userID string) (*ReadStatus, error)
	SaveReadStatus(*ReadStatus) (*ReadStatus, error)
}

// ReadStatus allows us to mark things read or unread
type ReadStatus struct {
	ID     string `json:"id"`
	PostID string `json:"post_id"`
	UserID string `json:"user_id"`

	ReadAt   time.Time `json:"read_at"`
	Device   string    `json:"device"`
	Location string    `json:"location"`
}
