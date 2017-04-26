package hydrocarbon

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// A UserStore is an interface used to seperate the UserAPI from knowledge of the
// actual underlying database
type UserStore interface {
	CreateUser(ctx context.Context, email string) (string, error)
	CreateLoginToken(ctx context.Context, userID, userAgent, ip string) (string, error)
	ActivateLoginToken(ctx context.Context, token string) (string, error)
	CreateSession(ctx context.Context, userID, userAgent, ip string) (string, error)
	DeactivateSession(ctx context.Context, key string) error
}

// UserAPI encapsulates everything related to user management
type UserAPI struct {
	s UserStore
	m Mailer
}

func NewUserAPI(s UserStore, m Mailer) *UserAPI {
	return &UserAPI{
		s: s,
		m: m,
	}
}

var (
	registerSuccess = []byte(`{"status":"success", "note": "check your email, token expires in 24 hours"}`)
)

func (ua *UserAPI) RequestToken(w http.ResponseWriter, r *http.Request) {
	var registerData struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(io.LimitReader(r.Body, 4*1024)).Decode(&registerData)
	if err != nil {
		panic(err)
		// do something
	}

	if len(registerData.Email) > 256 || !strings.Contains(registerData.Email, "@") {
		panic(err)
		// error out
	}

	userID, err := ua.s.CreateUser(r.Context(), registerData.Email)
	if err != nil {
		panic(err)
		// something
	}

	lt, err := ua.s.CreateLoginToken(r.Context(), userID, r.UserAgent(), r.RemoteAddr)
	if err != nil {
		panic(err)
		// something
	}

	err = ua.m.Send(registerData.Email, "go to $URL with "+lt)
	if err != nil {
		panic(err)
		// something
	}

	w.Write(registerSuccess)
}

func (ua *UserAPI) Activate(w http.ResponseWriter, r *http.Request) {
	var activateData struct {
		Token string `json:"token"`
	}

	err := json.NewDecoder(io.LimitReader(r.Body, 4*1024)).Decode(&activateData)
	if err != nil {
		panic(err)
		// do something
	}

	userID, err := ua.s.ActivateLoginToken(r.Context(), activateData.Token)
	if err != nil {
		panic(err)
		// do something
	}

	key, err := ua.s.CreateSession(r.Context(), userID, r.UserAgent(), r.RemoteAddr)
	if err != nil {
		panic(err)
		// do something
	}

	var activateSuccess = struct {
		Status string `json:"status"`
		Key    string `json:"key"`
	}{
		"success",
		key,
	}

	err = json.NewEncoder(w).Encode(&activateSuccess)
	if err != nil {
		// do something
	}
}

func (ua *UserAPI) Deactivate(w http.ResponseWriter, r *http.Request) {
	var deactivateData struct {
		Key string `json:"key"`
	}

	err := json.NewDecoder(io.LimitReader(r.Body, 4*1024)).Decode(&deactivateData)
	if err != nil {
		panic(err)
		// do something
	}

	err = ua.s.DeactivateSession(r.Context(), deactivateData.Key)
	if err != nil {
		panic(err)
		// do something
	}

	var deactivateSuccess = struct {
		Status string `json:"status"`
	}{
		"success",
	}

	err = json.NewEncoder(w).Encode(&deactivateSuccess)
	if err != nil {
		// do something
	}
}
