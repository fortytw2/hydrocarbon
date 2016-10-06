package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fortytw2/kiasu"
	"github.com/fortytw2/kiasu/stores/mem"
)

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	u := mem.NewStore()
	m := kiasu.FakeMailer()

	token, err := u.CreateUser(context.Background(), m, "luke@jedicouncil.gov", "IamABest91030!")
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Errorf("no token! %s", token)
	}

	req, err := http.NewRequest("GET", "http://kiasu.io/auth_test", nil)
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	Authenticate(u, func(w http.ResponseWriter, r *http.Request) error {
		if activeUser, ok := r.Context().Value("user").(*kiasu.User); ok {
			if activeUser.Email != "luke@jedicouncil.gov" {
				t.Error("active user isn't luke!")

				return errors.New("something is wrong")
			}
		} else {
			t.Error("no active user in context!")
		}

		if accessTok, ok := r.Context().Value("access_token").(string); ok {
			if accessTok != token {
				t.Error("access_token in context != accessToken from store")
				return errors.New("something is wrong")
			}
		} else {
			t.Error("no access token in context!")
		}

		w.WriteHeader(http.StatusOK)
		return nil
	})(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("could not authenticate w/ token, status %d", w.Code)
	}
}
