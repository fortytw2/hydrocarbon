package mem

import (
	"errors"
	"sync"
	"time"

	"context"

	"github.com/fortytw2/abdi"
	"github.com/fortytw2/kiasu"
)

type store struct {
	signKey []byte

	users     []*kiasu.User
	userIDMax int
	userMu    sync.RWMutex

	sessions  []*kiasu.Session
	sessIDMax int
	sessMu    sync.RWMutex
}

// NewStore returns a purely memory backed store :)
func NewStore(signingKey []byte) kiasu.Store {
	return &store{
		signKey:  signingKey,
		users:    make([]*kiasu.User, 0),
		sessions: make([]*kiasu.Session, 0),
	}
}

func (s *store) GetUser(_ context.Context, accessToken string) (*kiasu.User, error) {
	s.sessMu.RLock()
	defer s.sessMu.RUnlock()
	for _, sesh := range s.sessions {
		if sesh.Token == accessToken {
			return s.userByID(sesh.ID)
		}
	}
	return nil, errors.New("invalid accesstoken or no user found")
}

func (s *store) userByID(id int) (*kiasu.User, error) {
	s.userMu.RLock()
	defer s.userMu.RUnlock()
	for _, u := range s.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("no user found")
}

func (s *store) CreateUser(_ context.Context, m kiasu.Mailer, email string, pw string) (string, error) {
	s.userMu.Lock()
	defer s.userMu.Unlock()
	s.userIDMax++
	h, err := abdi.Hash(pw, s.signKey)
	if err != nil {
		return "", err
	}
	t := randToken(24)
	s.users = append(s.users, &kiasu.User{
		ID:                s.userIDMax,
		Email:             email,
		EncryptedPassword: *h,
		ConfirmationToken: &t,
		Confirmed:         false,
		NotifyWindow:      30 * time.Second,
	})

	return t, nil
}

func (s *store) ActivateUser(_ context.Context, confirmToken string) (string, error) {
	s.userMu.RLock()
	for _, u := range s.users {
		if confirmToken == *u.ConfirmationToken {
			s.userMu.RUnlock()
			s.userMu.Lock()
			u.Confirmed = true
			u.Active = true
			s.userMu.Unlock()

			s.sessMu.Lock()
			s.sessIDMax++
			t := randToken(24)
			s.sessions = append(s.sessions, &kiasu.Session{
				ID:        s.sessIDMax,
				UserID:    u.ID,
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(time.Hour),
				Token:     t,
			})
			s.sessMu.Unlock()
			return t, nil
		}
	}
	s.userMu.RUnlock()
	return "", errors.New("could not activate")

}

func (s *store) NewSession(_ context.Context, email string, pw string) (string, error) {
	panic("not implemented")
}

func (s *store) GetActiveSessions(_ context.Context, accessToken string, p *kiasu.Pagination) ([]kiasu.Session, error) {
	u, err := s.GetUser(context.Background(), accessToken)
	if err != nil {
		return nil, err
	}

	var out []kiasu.Session
	s.sessMu.RLock()
	for _, s := range s.sessions {
		if s.UserID == u.ID {
			out = append(out, *s)
		}
	}
	s.sessMu.RUnlock()

	return out, nil
}

func (s *store) GetPastSessions(_ context.Context, accessToken string, p *kiasu.Pagination) ([]kiasu.Session, error) {
	panic("not implemented")
}

func (s *store) InvalidateToken(_ context.Context, accessToken string) error {
	panic("not implemented")
}

func (s *store) GetFeeds(_ context.Context, p *kiasu.Pagination) ([]kiasu.Feed, error) {
	panic("not implemented")
}

func (s *store) GetUserFeeds(_ context.Context, accessToken string, p *kiasu.Pagination) ([]kiasu.Feed, error) {
	panic("not implemented")
}

func (s *store) ReOrderFeed(_ context.Context, accessToken string, feedID string, newOrder int) ([]kiasu.Feed, error) {
	panic("not implemented")
}

func (s *store) GetFeedPosts(_ context.Context, accessToken string, feedID string, p *kiasu.Pagination) ([]kiasu.Post, error) {
	panic("not implemented")
}

func (s *store) Charge(_ context.Context, accessToken string, chargeToken string) error {
	panic("not implemented")
}

func (s *store) GetUsersByExpiry(_ context.Context, m kiasu.Mailer, expireAfter time.Time, p *kiasu.Pagination) ([]kiasu.User, error) {
	panic("not implemented")
}

func (s *store) AddSubscription(_ context.Context, email string, activeUntil time.Time) error {
	panic("not implemented")
}
