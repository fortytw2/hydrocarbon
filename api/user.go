package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fortytw2/kiasu"
	"github.com/go-kit/kit/log"
)

// UserProfile encodes the current user's profile to JSON
func UserProfile(l log.Logger, us kiasu.UserStore) ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		user := r.Context().Value("user").(*kiasu.User)
		if err := activeUser(user); err != nil {
			return err
		}

		err := json.NewEncoder(w).Encode(user)
		if err != nil {
			return NewHTTPError("could not write response", http.StatusInternalServerError)
		}

		return nil
	}
}

// UserSessions lists all of the active users sessions
func UserSessions(l log.Logger, us kiasu.UserStore) ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		user := r.Context().Value("user").(*kiasu.User)
		if err := activeUser(user); err != nil {
			return err
		}

		return nil
	}
}

// ConfirmToken confirms an authentication token, activating a user
func ConfirmToken(l log.Logger, us kiasu.UserStore) ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		token := r.URL.Query().Get("token")
		if token == "" {
			return NewHTTPError("not a valid confirmation token", http.StatusBadRequest)
		}

		accessToken, err := us.ActivateUser(r.Context(), token)
		if err != nil {
			return NewHTTPError("could not activate user", http.StatusBadRequest)
		}

		fmt.Fprintf(w, `{"access_token": "%s"}`, accessToken)

		return nil
	}
}

// RegisterUser creates a new user
func RegisterUser(l log.Logger, m kiasu.Mailer, us kiasu.UserStore) ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			return NewHTTPError("error parsing json", http.StatusBadRequest)
		}

		sesh, err := us.CreateUser(r.Context(), m, body.Email, body.Password)
		if err != nil {
			return NewHTTPError("could not register user", http.StatusForbidden)
		}

		// this needs to sent via email
		fmt.Fprintf(w, `{"confirmation_token": "%s"}`, sesh)
		w.WriteHeader(http.StatusOK)
		return nil
	}
}

// Login logs a user in, giving them a fresh access token
func Login(l log.Logger, us kiasu.UserStore) ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			return NewHTTPError("error parsing json", http.StatusBadRequest)
		}

		sesh, err := us.NewSession(r.Context(), body.Email, body.Password)
		if err != nil {
			return NewHTTPError("could not create a new session", http.StatusInternalServerError)
		}

		fmt.Fprintf(w, `{"access_token": "%s"}`, sesh)
		w.WriteHeader(http.StatusOK)
		return nil
	}
}

// Logout invalidates the current token
func Logout(l log.Logger, us kiasu.UserStore) ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		token := r.Context().Value("access_token").(string)
		err := us.InvalidateToken(r.Context(), token)
		if err != nil {
			return NewHTTPError("could not invalidate current session", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		return nil
	}
}

// DeactivateUser deletes the users account
func DeactivateUser(l log.Logger, us kiasu.UserStore) ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}

func activeUser(u *kiasu.User) error {
	if !u.Active {
		return NewHTTPError("inactive or deleted user account", http.StatusUnauthorized)
	}

	return nil
}
