package sentry

import (
	"context"
	"net/http"

	"github.com/fortytw2/hydrocarbon/discollect"
)

// ErrorReporter implements discollect.ErrorReporter and hydrocarbon.ErrorReporter
// using Sentry.io
type ErrorReporter struct {
	pubKey, privKey string
	c               *http.Client

	shutdown chan chan struct{}
	reports  chan *reportRequest
}

// NewErrorReporter instantiates a new ErrorReporter
func NewErrorReporter(c *http.Client, pubKey, privKey string) (*ErrorReporter, error) {
	return &ErrorReporter{
		pubKey:  pubKey,
		privKey: privKey,
		c:       c,

		shutdown: make(chan chan struct{}),
		reports:  make(chan *reportRequest, 64),
	}, nil
}

// Start launches the ErrorReporter
func (er *ErrorReporter) Start() {
	for {
		select {
		case a := <-er.shutdown:
			a <- struct{}{}
			return
			// case r := <-er.reports:
			// er.sendToSentry(r)
		}
	}
}

func (er *ErrorReporter) Stop() {
	c := make(chan struct{})
	er.shutdown <- c
	<-c
}

type reportRequest struct {
	ro  *discollect.ReporterOpts
	err error
}

// Report enqueues an error to be sent to sentry in the future
func (er *ErrorReporter) Report(ctx context.Context, ro *discollect.ReporterOpts, err error) {
	er.reports <- &reportRequest{
		ro:  ro,
		err: err,
	}
}

// func (er *ErrorReporter) sendToSentry(r *reportRequest) {
// 	xSentryAuth := fmt.Sprintf(`Sentry sentry_version=9,
// sentry_client=hydrocarbon/0.1,
// sentry_timestamp=%d,
// sentry_key=%s,
// sentry_secret=%s`, time.Now().UnixNano(), er.pubKey, er.privKey)

// 	req := http.NewRequest("POST", "https://sentry.io/api/")

// }
