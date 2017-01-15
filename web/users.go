package web

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/internal/httputil"
	"github.com/fortytw2/hydrocarbon/internal/token"
	"github.com/mholt/binding"
)

// login renders a dummy page for logging in
func renderLogin(w http.ResponseWriter, r *http.Request) error {
	out, err := TMPLERRlogin("Hydrocarbon", false, 0)
	if err != nil {
		return httputil.Wrap(err, http.StatusInternalServerError)
	}

	_, err = w.Write([]byte(out))
	if err != nil {
		return httputil.Wrap(err, http.StatusInternalServerError)
	}

	return nil
}

// register displays a sign up page
func renderRegister(w http.ResponseWriter, r *http.Request) error {
	out := TMPLregister("Hydrocarbon", false, 0)
	_, err := w.Write([]byte(out))
	if err != nil {
		return httputil.Wrap(err, http.StatusInternalServerError)
	}

	return nil
}

type registration struct {
	Email    string
	Password string
}

// Then provide a field mapping (pointer receiver is vital)
func (r *registration) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&r.Email:    "email",
		&r.Password: "password",
	}
}

// newUser processes a post request
func newUser(s *hydrocarbon.Store) httputil.ErrorHandler {
	return func(w http.ResponseWriter, req *http.Request) error {
		r := new(registration)
		errs := binding.Bind(req, r)
		if errs.Handle(w) {
			return nil
		}

		user, err := s.CreateUser(r.Email, r.Password)
		if err != nil {
			return err
		}

		newToken, err := token.GenerateRandomString(32)
		if err != nil {
			return err
		}
		// right here should send a confirmation email
		// user.ConfirmationToken
		sesh, err := s.Sessions.CreateSession(&hydrocarbon.Session{
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(14 * 24 * time.Hour),
			Token:     newToken,
		})
		if err != nil {
			return err
		}

		js, _ := json.Marshal(sesh)
		_, err = w.Write(js)

		return err
	}
}

// confirmUser asserts that the user has a valid email
func confirmUser(w http.ResponseWriter, r *http.Request) {

}

// forgotPassword sends a reset email
func forgotPassword(w http.ResponseWriter, r *http.Request) {

}

// deleteSession invalidates an existing session
func deleteSession(w http.ResponseWriter, r *http.Request) {

}

// newSession creates a new session
func newSession(w http.ResponseWriter, r *http.Request) {

}
