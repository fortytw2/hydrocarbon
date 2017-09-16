package hydrocarbon

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

// A UserStore is an interface used to seperate the UserAPI from knowledge of the
// actual underlying database
type UserStore interface {
	CreateOrGetUser(ctx context.Context, email string) (string, error)
	CreateLoginToken(ctx context.Context, userID, userAgent, ip string) (string, error)
	ActivateLoginToken(ctx context.Context, token string) (string, error)
	CreateSession(ctx context.Context, userID, userAgent, ip string) (string, string, error)
	ListSessions(ctx context.Context, key string, page int) ([]*Session, error)
	DeactivateSession(ctx context.Context, key string) error
}

// UserAPI encapsulates everything related to user management
type UserAPI struct {
	s  UserStore
	m  Mailer
	ks *KeySigner
}

func NewUserAPI(s UserStore, ks *KeySigner, m Mailer) *UserAPI {
	return &UserAPI{
		s:  s,
		ks: ks,
		m:  m,
	}
}

var (
	registerSuccess = []byte(`{"status":"success", "note": "email sent, token expires in 24 hours"}`)
)

func (ua *UserAPI) RequestToken(w http.ResponseWriter, r *http.Request) {
	var registerData struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(io.LimitReader(r.Body, 4*1024)).Decode(&registerData)
	if err != nil {
		writeErr(w, err)
		return
	}

	if len(registerData.Email) > 256 || !strings.Contains(registerData.Email, "@") {
		writeErr(w, err)
		return
	}

	userID, err := ua.s.CreateOrGetUser(r.Context(), registerData.Email)
	if err != nil {
		writeErr(w, err)
		return
	}

	lt, err := ua.s.CreateLoginToken(r.Context(), userID, r.UserAgent(), getRemoteIP(r))
	if err != nil {
		writeErr(w, err)
		return
	}

	err = ua.m.Send(registerData.Email, fmt.Sprintf("visit %s/login-callback?token=%s to login", ua.m.RootDomain(), lt))
	if err != nil {
		writeErr(w, err)
		return
	}

	w.Write(registerSuccess)
}

// ListSessions writes out all of a users current / past sessions
func (ua *UserAPI) ListSessions(w http.ResponseWriter, r *http.Request) {
	key, err := ua.ks.Verify(r.Header.Get("X-Hydrocarbon-Key"))
	if err != nil {
		writeErr(w, err)
		return
	}

	sess, err := ua.s.ListSessions(r.Context(), key, 0)
	if err != nil {
		writeErr(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(sess)
	if err != nil {
		writeErr(w, err)
		return
	}
}

func (ua *UserAPI) Activate(w http.ResponseWriter, r *http.Request) {
	var activateData struct {
		Token string `json:"token"`
	}

	err := json.NewDecoder(io.LimitReader(r.Body, 4*1024)).Decode(&activateData)
	if err != nil {
		writeErr(w, err)
		return
	}

	userID, err := ua.s.ActivateLoginToken(r.Context(), activateData.Token)
	if err != nil {
		writeErr(w, err)
		return
	}

	email, key, err := ua.s.CreateSession(r.Context(), userID, r.UserAgent(), getRemoteIP(r))
	if err != nil {
		writeErr(w, err)
		return
	}

	key, err = ua.ks.Sign(key)
	if err != nil {
		writeErr(w, err)
		return
	}

	var activateSuccess = struct {
		Status string `json:"status"`
		Email  string `json:"email"`
		Key    string `json:"key"`
	}{
		"success",
		email,
		key,
	}

	err = json.NewEncoder(w).Encode(&activateSuccess)
	if err != nil {
		// do something
	}
}

func (ua *UserAPI) Deactivate(w http.ResponseWriter, r *http.Request) {
	key, err := ua.ks.Verify(r.Header.Get("X-Hydrocarbon-Key"))
	if err != nil {
		writeErr(w, err)
		return
	}

	err = ua.s.DeactivateSession(r.Context(), key)
	if err != nil {
		writeErr(w, err)
		return
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

func getRemoteIP(r *http.Request) string {
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	fwdIP := r.Header.Get("X-Forwarded-For")
	fwdSplit := strings.Split(fwdIP, ",")
	if fwdIP != "" {
		// pick the leftmost x-forwarded-for addr
		return fwdSplit[0]
	}

	// this literally can't fail on r.RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func writeErr(w http.ResponseWriter, err error) {
	var s = struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}{
		"error",
		err.Error(),
	}

	json.NewEncoder(w).Encode(s)
	w.WriteHeader(http.StatusOK)
}
