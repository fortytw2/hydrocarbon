package tmpl

import "github.com/fortytw2/hydrocarbon"

// Store provides basic persistence primitives
type Store struct{}

// NewStore creates a primitive persistence layer
func NewStore() (hydrocarbon.PrimitiveStore, error) {
	return &Store{}, nil
}
