package hydrocarbon

import "net/http"

// Doer is a simplified interface to allow more complex *http.Client
// implementations to be passed around
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}
