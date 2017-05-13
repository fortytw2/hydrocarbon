package postmark

import "log"

// Mailer sends mails via Postmark
type Mailer struct {
	Key    string
	Domain string
}

func (m *Mailer) Send(email string, body string) error {
	log.Println("hydrocarbon: new mail to", email, "\n", body)
	return nil
}

// RootDomain returns the StdoutMailer root domain
func (m *Mailer) RootDomain() string {
	return m.Domain
}
