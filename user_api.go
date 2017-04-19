package hydrocarbon

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// A UserStore is an interface used to seperate the UserAPI from knowledge of the
// actual underlying database
type UserStore interface {
	CreateUser(ctx context.Context, email string) (string, error)
	CreateLoginToken(ctx context.Context, userID string) (string, error)
	ActivateLoginToken(ctx context.Context, token string) (string, error)
	CreateSession(ctx context.Context, userID, userAgent, ip string) (string, error)
	DeleteSession(ctx context.Context, token string) error
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

var registerSuccess = []byte(`{"status":"success, check your email"}`)

func (ua *UserAPI) Register(w http.ResponseWriter, r *http.Request) {
	var registerData = struct {
		Email string `json:"email"`
	}{}

	err := json.NewDecoder(io.LimitReader(r.Body, 4*1024)).Decode(&registerData)
	if err != nil {
		// do something
	}

	if len(registerData.Email) > 256 {
		// error out
	}

	userID, err := ua.s.CreateUser(r.Context(), registerData.Email)
	if err != nil {
		// something
	}

	lt, err := ua.s.CreateLoginToken(r.Context(), userID)
	if err != nil {
		// something
	}

	err = ua.m.Send(registerData.Email, "go to $URL with "+lt)
	if err != nil {
		// something
	}

	w.Write(registerSuccess)
}
