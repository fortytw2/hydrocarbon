package hydrocarbon

import "time"

// A Folder is a collection of feeds
type Folder struct {
	ID string `json:"id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name  string `json:"name"`
	Feeds []Feed `json:"feeds"`
}
