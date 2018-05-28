package discollect

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/google/uuid"
)

// A Writer is able to process and output datums retrieved by Discollect plugins
type Writer interface {
	Write(ctx context.Context, scrapeID uuid.UUID, f interface{}) error
	io.Closer
}

// StdoutWriter fmt.Printfs to stdout
type StdoutWriter struct{}

// Write printf %+v the datum to stdout
func (sw *StdoutWriter) Write(ctx context.Context, _ uuid.UUID, f interface{}) error {
	return json.NewEncoder(os.Stdout).Encode(f)
}

// Close is a no-op function so the StdoutWriter works
func (sw *StdoutWriter) Close() error {
	return nil
}
