package api

import (
	"net/http"
	"strings"

	"github.com/fortytw2/kiasu"
	"golang.org/x/net/context"
)

// Authenticate wraps any given handler in basic authentication
// using Bearer tokens
func Authenticate(us kiasu.UserStore, x http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := w.Header().Get("Authorization")
		trimmed := strings.TrimLeft(authHeader, "Bearer ")
		if trimmed == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := us.GetUser(r.Context(), trimmed)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// add the user + access_token to context
		newCtx := context.WithValue(r.Context(), "user", user)
		r.WithContext(context.WithValue(newCtx, "access_token", trimmed))

		x(w, r)
	}
}
