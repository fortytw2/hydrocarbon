package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"context"

	"github.com/fortytw2/kiasu"
	"github.com/fortytw2/kiasu/stores/mem"
	"github.com/go-kit/kit/log"
)

func TestConfirmUser(t *testing.T) {
	t.Parallel()

	u := mem.NewStore()
	m := kiasu.FakeMailer()

	confirmToken, err := u.CreateUser(context.Background(), m, "luke@jedicouncil.gov", "IamABest91030!")
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("GET", "http://kiasu.io/api/v1/users/confirm?token="+confirmToken, nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	ConfirmToken(log.NewNopLogger(), u).ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("confirm token did not return 200 - %d", w.Code)
		return
	}

}

func TestUserProfile(t *testing.T) {
	t.Parallel()

	u := mem.NewStore()
	m := kiasu.FakeMailer()

	confirmToken, err := u.CreateUser(context.Background(), m, "luke@jedicouncil.gov", "IamABest91030!")
	if err != nil {
		t.Error(err)
	}

	accessToken, _ := u.ActivateUser(nil, confirmToken)

	req, err := http.NewRequest("GET", "http://kiasu.io/api/v1/users/", nil)
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w := httptest.NewRecorder()
	Authenticate(u, UserProfile(log.NewNopLogger(), u)).ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("user profile did not return 200 - %d", w.Code)
		return
	}

	var us kiasu.User
	err = json.NewDecoder(w.Body).Decode(&us)
	if err != nil {
		t.Error(err)
	}

	if us.Email != "luke@jedicouncil.gov" {
		t.Errorf("user email is NOT luke! - %s", us.Email)
	}
}

func TestRegisterUser(t *testing.T) {
	t.Parallel()

}
