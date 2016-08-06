package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fortytw2/kiasu"
	"github.com/go-kit/kit/log"
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	"golang.org/x/net/context"
)

// AddUserRoutes adds user routes under a *xmux.Mux
func AddUserRoutes(xm *xmux.Group, l log.Logger, m kiasu.Mailer, db kiasu.Store) {
	users := xm.NewGroup("/user")
	users.GET("/me", UserProfile(l, db))
	users.GET("/sessions", UserSessions(l, db))
	users.DELETE("/sessions", Logout(l, db))

	users.GET("/confirm", ConfirmToken(l, m, db))
	users.POST("/login", Login(l, m, db))
	users.POST("/new", RegisterUser(l, m, db))
	users.DELETE("/deactivate", DeactivateUser(l, db))
}

// UserProfile encodes the current user's profile to JSON
func UserProfile(l log.Logger, us kiasu.UserStore) xhandler.HandlerFuncC {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		user := ctx.Value("user").(*kiasu.User)
		err := json.NewEncoder(w).Encode(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// UserSessions lists all of the active users sessions
func UserSessions(l log.Logger, us kiasu.UserStore) xhandler.HandlerFuncC {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	}
}

// ConfirmToken confirms an authentication token, activating a user
func ConfirmToken(l log.Logger, m kiasu.Mailer, us kiasu.UserStore) xhandler.HandlerFuncC {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	}
}

// RegisterUser creates a new user
func RegisterUser(l log.Logger, m kiasu.Mailer, us kiasu.UserStore) xhandler.HandlerFuncC {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sesh, err := us.CreateUser(ctx, m, body.Email, body.Password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fmt.Fprintf(w, `{"access_token": "%s"}`, sesh)
		w.WriteHeader(http.StatusOK)
	}
}

// Login logs a user in, giving them a fresh access token
func Login(l log.Logger, m kiasu.Mailer, us kiasu.UserStore) xhandler.HandlerFuncC {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sesh, err := us.NewSession(ctx, m, body.Email, body.Password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fmt.Fprintf(w, `{"access_token": "%s"}`, sesh)
		w.WriteHeader(http.StatusOK)
	}
}

// Logout invalidates the current token
func Logout(l log.Logger, us kiasu.UserStore) xhandler.HandlerFuncC {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		token := ctx.Value("access_token").(string)
		err := us.InvalidateToken(ctx, token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// DeactivateUser deletes the users account
func DeactivateUser(l log.Logger, us kiasu.UserStore) xhandler.HandlerFuncC {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	}
}
