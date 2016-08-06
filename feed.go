package kiasu

import "time"

// A Feed is a single encapsulating unit around a source of news / posts / whatever
type Feed struct {
	ID       int  `json:"id"`
	PluginID *int `json:"plugin_id,omitempty"`

	CreatedAt   time.Time  `json:"created_at"`
	RefreshedAt *time.Time `json:"refreshed_at"`

	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`

	HexColor string `json:"hex_color"`
	IconURL  string `json:"icon_url"`
}
