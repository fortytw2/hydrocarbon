package postmark

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fortytw2/hydrocarbon/httpx"
)

// Mailer sends mails via Postmark
type Mailer struct {
	Key    string
	Domain string
	Client *http.Client
}

type mailReq struct {
	From     string `json:"From"`
	To       string `json:"To"`
	Cc       string `json:"Cc"`
	Bcc      string `json:"Bcc"`
	Subject  string `json:"Subject"`
	Tag      string `json:"Tag"`
	HTMLBody string `json:"HtmlBody"`
	TextBody string `json:"TextBody"`
	ReplyTo  string `json:"ReplyTo"`
	Headers  []struct {
		Name  string `json:"Name"`
		Value string `json:"Value"`
	} `json:"Headers"`
	TrackOpens bool   `json:"TrackOpens"`
	TrackLinks string `json:"TrackLinks"`
}

// Send sends a mail using the postmark api
func (m *Mailer) Send(email, subject, body string) error {
	buf, err := json.Marshal(&mailReq{
		From:     "support@hydrocarbon.io",
		To:       email,
		Subject:  subject,
		HTMLBody: body,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.postmarkapp.com/email", bytes.NewReader(buf))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", m.Key)

	resp, err := m.Client.Do(req)
	if err != nil {
		return err
	}

	err = httpx.DrainAndClose(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusUnprocessableEntity {
		return errors.New("error sending to postmark, got 422")
	}

	return nil
}

// RootDomain returns the StdoutMailer root domain
func (m *Mailer) RootDomain() string {
	return m.Domain
}
