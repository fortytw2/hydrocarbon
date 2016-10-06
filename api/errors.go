package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type httpError struct {
	Message string `json:"message"`
	Code    int    `json:"-"`
}

// NewHTTPError returns a custom error type with a response code attached
// that cleanly marshals to JSON
func NewHTTPError(msg string, code int) error {
	return httpError{
		Message: msg,
		Code:    code,
	}
}

func (h httpError) Error() string {
	return h.Message
}

func (h httpError) Write(w http.ResponseWriter) {
	w.WriteHeader(h.Code)
	json.NewEncoder(w).Encode(h)
}

// ErrorHandler is an HTTP Handler mixed with a error to avoid repetitive
// writeHeader ; return
type ErrorHandler func(w http.ResponseWriter, r *http.Request) error

func (eh ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancelFunc := context.WithDeadline(r.Context(), time.Now().Add(time.Second*10))
	defer cancelFunc()

	r.WithContext(ctx)
	err := eh(w, r)
	if err != nil {
		if he, ok := err.(httpError); ok {
			he.Write(w)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
