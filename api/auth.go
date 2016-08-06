package api

import (
	"net/http"
	"strings"

	"github.com/fortytw2/kiasu"
	"github.com/rs/xhandler"
	"golang.org/x/net/context"
)

// Authenticate wraps any given handler in basic authentication
// using Bearer tokens
func Authenticate(us kiasu.UserStore, x xhandler.HandlerFuncC) xhandler.HandlerFuncC {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		authHeader := w.Header().Get("Authorization")
		trimmed := strings.TrimLeft(authHeader, "Bearer ")
		if trimmed == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := us.GetUser(ctx, trimmed)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		newCtx := context.WithValue(ctx, "user", user)
		x(context.WithValue(newCtx, "access_token", trimmed), w, r)
	}
}
