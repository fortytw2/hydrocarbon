package hydrocarbon

import "time"

// FolderStore saves and persists folders
type FolderStore interface {
	CreateFolder(*Folder) (*Folder, error)
	AddFeed(folderID, feedID string) error
}

// A Folder is a collection of feeds
type Folder struct {
	ID string `json:"id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name  string `json:"name"`
	Feeds []Feed `json:"feeds"`
}
