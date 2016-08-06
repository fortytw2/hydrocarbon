package kiasu

// Pagination is used to pass pagination information around, as well as
// provide helper methods (DRY)
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}
