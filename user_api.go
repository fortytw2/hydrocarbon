package hydrocarbon

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
)

// A UserStore is an interface used to seperate the UserAPI from knowledge of the
// actual underlying database
type UserStore interface {
	CreateOrGetUser(ctx context.Context, email string) (string, bool, error)
	SetStripeIDs(ctx context.Context, userID, customerID, subscriptionID string) error

	CreateLoginToken(ctx context.Context, userID, userAgent, ip string) (string, error)
	ActivateLoginToken(ctx context.Context, token string) (string, error)

	CreateSession(ctx context.Context, userID, userAgent, ip string) (string, string, error)
	ListSessions(ctx context.Context, key string, page int) ([]*Session, error)
	DeactivateSession(ctx context.Context, key string) error
}

// UserAPI encapsulates everything related to user management
type UserAPI struct {
	paymentRequired bool
	stripePlanID    string
	sc              *client.API
	s               UserStore
	m               Mailer
	ks              *KeySigner
}

// NewUserAPI sets up a new UserAPI used for user/session management
func NewUserAPI(s UserStore, ks *KeySigner, m Mailer, stripePlanID, stripeKey string, paymentRequired bool) *UserAPI {
	var c *client.API
	if paymentRequired {
		c = &client.API{}
		c.Init(stripeKey, nil)
	}

	return &UserAPI{
		s:               s,
		ks:              ks,
		m:               m,
		sc:              c,
		stripePlanID:    stripePlanID,
		paymentRequired: paymentRequired,
	}
}

// RequestToken emails a token that can be exchanged for a session
func (ua *UserAPI) RequestToken(w http.ResponseWriter, r *http.Request) error {
	var registerData struct {
		Email string `json:"email"`
	}

	err := limitDecoder(r, &registerData)
	if err != nil {
		return err
	}

	if len(registerData.Email) == 0 || len(registerData.Email) > 128 || !strings.Contains(registerData.Email, "@") {
		return errors.New("invalid email")
	}

	userID, paid, err := ua.s.CreateOrGetUser(r.Context(), registerData.Email)
	if err != nil {
		return err
	}

	if ua.paymentRequired && !paid {
		return errors.New("payment is required")
	}

	lt, err := ua.s.CreateLoginToken(r.Context(), userID, r.UserAgent(), GetRemoteIP(r))
	if err != nil {
		return err
	}

	err = ua.m.Send(registerData.Email, fmt.Sprintf("visit %s/callback?token=%s to login", ua.m.RootDomain(), lt))
	if err != nil {
		return err
	}

	return writeSuccess(w, "check your email for a login token, token expires in 24 hours")
}

// CreatePayment sets up the initial stripe stuff for a user
func (ua *UserAPI) CreatePayment(w http.ResponseWriter, r *http.Request) error {
	if !ua.paymentRequired {
		return errors.New("payments are not enabled on this instance")
	}

	var stripeData struct {
		Email  string `json:"email"`
		Coupon string `json:"coupon"`
		Token  string `json:"token"`
	}

	err := limitDecoder(r, &stripeData)
	if err != nil {
		return err
	}

	userID, paid, err := ua.s.CreateOrGetUser(r.Context(), stripeData.Email)
	if err != nil {
		return err
	}

	if paid {
		return errors.New("subscription already exists")
	}

	params := &stripe.CustomerParams{
		Email: &stripeData.Email,
	}
	err = params.SetSource(stripeData.Token)
	if err != nil {
		return err
	}

	customer, err := customer.New(params)
	if err != nil {
		return err
	}

	sp := &stripe.SubscriptionParams{
		Customer: &customer.ID,
		Plan:     &ua.stripePlanID,
	}

	if stripeData.Coupon != "" {
		sp.Coupon = &stripeData.Coupon
	}

	s, err := sub.New(sp)
	if err != nil {
		return err
	}

	err = ua.s.SetStripeIDs(r.Context(), userID, customer.ID, s.ID)
	if err != nil {
		return err
	}

	return writeSuccess(w, "stripe subscription created")
}

// ListSessions writes out all of a users current / past sessions
func (ua *UserAPI) ListSessions(w http.ResponseWriter, r *http.Request) error {
	key, err := ua.ks.Verify(r.Header.Get("X-Hydrocarbon-Key"))
	if err != nil {
		return err
	}

	sess, err := ua.s.ListSessions(r.Context(), key, 0)
	if err != nil {
		return err
	}

	return writeSuccess(w, sess)
}

// Activate exchanges a token for a session key that can be used to make
// authenticated requests
func (ua *UserAPI) Activate(w http.ResponseWriter, r *http.Request) error {
	var activateData struct {
		Token string `json:"token"`
	}

	err := limitDecoder(r, &activateData)
	if err != nil {
		return err
	}

	userID, err := ua.s.ActivateLoginToken(r.Context(), activateData.Token)
	if err != nil {
		return err
	}

	email, key, err := ua.s.CreateSession(r.Context(), userID, r.UserAgent(), GetRemoteIP(r))
	if err != nil {
		return err
	}

	key, err = ua.ks.Sign(key)
	if err != nil {
		return err
	}

	var activationData = struct {
		Email string `json:"email"`
		Key   string `json:"key"`
	}{
		email,
		key,
	}

	return writeSuccess(w, activationData)
}

// Deactivate disables a key that the user is currently using
func (ua *UserAPI) Deactivate(w http.ResponseWriter, r *http.Request) error {
	key, err := ua.ks.Verify(r.Header.Get("X-Hydrocarbon-Key"))
	if err != nil {
		return err
	}

	err = ua.s.DeactivateSession(r.Context(), key)
	if err != nil {
		return err
	}

	return writeSuccess(w, nil)
}

func GetRemoteIP(r *http.Request) string {
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
