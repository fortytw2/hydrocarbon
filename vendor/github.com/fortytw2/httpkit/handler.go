package httpkit

import (
	"encoding/json"
	"net/http"
)

// An ErrorHandler is a HTTP Handler that returns an error type
type ErrorHandler func(w http.ResponseWriter, r *http.Request) error

func (e ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := e(w, r)
	if err != nil {
		httpErr, ok := err.(Error)
		if !ok {
			w.WriteHeader(500)
			return
		}

		err := json.NewEncoder(w).Encode(map[string]string{"error": httpErr.Error()})
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(httpErr.Status())
	}
}
