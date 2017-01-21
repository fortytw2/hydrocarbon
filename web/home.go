package web

import "net/http"

func renderHome(w http.ResponseWriter, r *http.Request) error {
	out, err := TMPLERRhome("Hydrocarbon", loggedIn(r))
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(out))
	return err
}

func renderPrivacy(w http.ResponseWriter, r *http.Request) error {
	out, err := TMPLERRprivacy("Hydrocarbon", loggedIn(r))
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(out))
	return err
}
