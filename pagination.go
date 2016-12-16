package hydrocarbon

// Pagination is used to pass pagination information around, as well as
// provide helper methods (DRY)
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`

	// SortBy is the field to sort by, order is asc or desc
	SortBy string `json:"sort_by"`
	Order  string `json:"order"`
}
