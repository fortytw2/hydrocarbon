package discollect

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrRateLimitExceeded is thrown when the rate limit is exceeded
	ErrRateLimitExceeded = errors.New("discollect: rate limit exceeded")
)

// RateLimit is a wrapper struct around a variety of per-config rate limits
type RateLimit struct {
	// Rate a single IP can make requests per second
	PerIP float64
	// Rate the entire scrape can operate at per second
	PerScrape float64
	// Rate per domain using the publicsuffix list to differentiate per second
	PerDomain float64
}

// A Limiter is used for per-site and per-config rate limits
// abstracted out into an interface so that distributed rate limiting
// is practical
type Limiter interface {
	// ReserveN returns a Reservation that indicates how long the caller must
	// wait before n events happen. The Limiter takes this Reservation into
	// account when allowing future events. ReserveN returns false if n exceeds
	// the Limiter's burst size.
	Reserve(rl *RateLimit, url string, scrapeID uuid.UUID) (Reservation, error)
}

// A Reservation holds information about events that are permitted by a Limiter
// to happen after a delay. A Reservation may be canceled, which may enable the
// Limiter to permit additional events.
type Reservation interface {
	// Cancel indicates that the reservation holder will not perform the
	// reserved action and reverses the effects of this Reservation on the rate
	// limit as much as possible, considering that other reservations may have
	// already been made.
	Cancel()
	// OK returns whether the limiter can provide the requested number of tokens
	// within the maximum wait time. If OK is false, Delay returns InfDuration,
	// and Cancel does nothing.
	OK() bool
	// Delay returns the duration for which the reservation holder must wait
	// before taking the reserved action. Zero duration means act immediately.
	// InfDuration means the limiter cannot grant the tokens requested in this
	// Reservation within the maximum wait time.
	Delay() time.Duration
}

// A NilLimiter is a Limiter that doesn't restrict anything
type NilLimiter struct{}

// Reserve returns a dummy reservation that always waits one second
func (*NilLimiter) Reserve(rl *RateLimit, url string, scrapeID uuid.UUID) (Reservation, error) {
	return &nilReservation{}, nil
}

type nilReservation struct{}

func (*nilReservation) Cancel() {}

func (*nilReservation) OK() bool {
	return true
}

func (*nilReservation) Delay() time.Duration {
	return time.Second
}
