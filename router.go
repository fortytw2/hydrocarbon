package hydrocarbon

import "net/http"

func NewRouter(ua *UserAPI) *http.ServeMux {
	m := http.NewServeMux()

	m.HandleFunc("/api/request", ua.Register)
	m.HandleFunc("/api/activate", ua.Activate)

	return m
}
