package kiasu

// A Folder is used to nest feeds for easy browsing. Only one level deep though
type Folder struct {
	ID int `json:"-"`

	Name  string `json:"name"`
	Feeds []Feed `json:"feeds"`
}
