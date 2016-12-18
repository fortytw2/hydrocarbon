package web

import "net/http"

func renderHome(w http.ResponseWriter, r *http.Request) error {
	out, err := TMPLERRhome("Hydrocarbon", false, 0)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(out))
	return err
}
