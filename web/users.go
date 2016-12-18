package web

import "net/http"

// login renders a dummy page for logging in
func renderLogin(w http.ResponseWriter, r *http.Request) {
	out, err := TMPLERRlogin("Hydrocarbon", false, 0)
	if err != nil {
		panic(err)
	}

	_, err = w.Write([]byte(out))
	if err != nil {
		panic(err)
	}
}

// register displays a sign up page
func renderRegister(w http.ResponseWriter, r *http.Request) {
	out, err := TMPLERRregister("Hydrocarbon", false, 0)
	if err != nil {
		panic(err)
	}

	_, err = w.Write([]byte(out))
	if err != nil {
		panic(err)
	}
}

// newUser processes a post request
func newUser(w http.ResponseWriter, r *http.Request) {

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
