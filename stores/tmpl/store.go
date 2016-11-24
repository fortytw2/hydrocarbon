package tmpl

import "github.com/fortytw2/kiasu"

// Store provides basic persistence primitives
type Store struct{}

// NewStore creates a primitive persistence layer
func NewStore() (kiasu.PrimitiveStore, error) {
	return &Store{}, nil
}
