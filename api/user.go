package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fortytw2/kiasu"
	"github.com/go-kit/kit/log"
)

// UserProfile encodes the current user's profile to JSON
func UserProfile(l log.Logger, us kiasu.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*kiasu.User)
		err := json.NewEncoder(w).Encode(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// UserSessions lists all of the active users sessions
func UserSessions(l log.Logger, us kiasu.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// ConfirmToken confirms an authentication token, activating a user
func ConfirmToken(l log.Logger, m kiasu.Mailer, us kiasu.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// RegisterUser creates a new user
func RegisterUser(l log.Logger, m kiasu.Mailer, us kiasu.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sesh, err := us.CreateUser(r.Context(), m, body.Email, body.Password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fmt.Fprintf(w, `{"access_token": "%s"}`, sesh)
		w.WriteHeader(http.StatusOK)
	}
}

// Login logs a user in, giving them a fresh access token
func Login(l log.Logger, m kiasu.Mailer, us kiasu.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sesh, err := us.NewSession(r.Context(), m, body.Email, body.Password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fmt.Fprintf(w, `{"access_token": "%s"}`, sesh)
		w.WriteHeader(http.StatusOK)
	}
}

// Logout invalidates the current token
func Logout(l log.Logger, us kiasu.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Context().Value("access_token").(string)
		err := us.InvalidateToken(r.Context(), token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// DeactivateUser deletes the users account
func DeactivateUser(l log.Logger, us kiasu.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
