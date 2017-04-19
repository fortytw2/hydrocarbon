package hydrocarbon

import "net/http"

func NewRouter(ua *UserAPI) *http.ServeMux {
	m := http.NewServeMux()

	// session management
	m.HandleFunc("/api/token/request", ua.RequestToken)
	m.HandleFunc("/api/token/activate", ua.Activate)
	m.HandleFunc("/api/key/deactivate", ua.Deactivate)

	return m
}
