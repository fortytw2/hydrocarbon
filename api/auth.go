package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/fortytw2/kiasu"
)

// Authenticate wraps any given handler in basic authentication
// using Bearer tokens
func Authenticate(us kiasu.UserStore, x ErrorHandler) ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		authHeader := r.Header.Get("Authorization")
		trimmed := strings.TrimLeft(authHeader, "Bearer ")
		if trimmed == "" {
			return NewHTTPError("invalid authorization header", http.StatusUnauthorized)
		}

		user, err := us.GetUser(r.Context(), trimmed)
		if err != nil {
			return NewHTTPError("invalid authorization header", http.StatusUnauthorized)
		}

		// add the user + access_token to context
		newCtx := context.WithValue(r.Context(), "user", user)
		newR := r.WithContext(context.WithValue(newCtx, "access_token", trimmed))

		return x(w, newR)
	}
}
