package hydrocarbon

import (
	"context"
	"log"
)

// ErrorReporter is used to report errors from long running / background jobs
// or to forward truly unknown errors to a service like Sentry
type ErrorReporter interface {
	Report(ctx context.Context, err error)
}

// StdoutReporter writes errors to Stdout
type StdoutReporter struct{}

// Report writes errors to stdout
func (s *StdoutReporter) Report(ctx context.Context, err error) {
	log.Println("error:", err)
}
