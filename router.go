package hydrocarbon

import "net/http"

func NewRouter(ua *UserAPI) *http.ServeMux {
	m := http.NewServeMux()

	m.HandleFunc("/api/register", ua.Register)

	return m
}
