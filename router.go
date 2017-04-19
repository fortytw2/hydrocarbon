package hydrocarbon

import "net/http"

func NewRouter(ua *UserAPI) *http.ServeMux {
	m := http.NewServeMux()

	// user management
	m.HandleFunc("/api/request", ua.RequestToken)
	m.HandleFunc("/api/activate", ua.Activate)
	m.HandleFunc("/api/deactivate", ua.Deactivate)

	return m
}
