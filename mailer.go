package hydrocarbon

import "log"

// A Mailer sends mail
type Mailer interface {
	Send(email string, body string) error
	RootDomain() string
}

// MockMailer is a fake mailer that records all mails sent
type MockMailer struct {
	Mails []string
}

// Send stores a mail in the local MockMailer
func (mm *MockMailer) Send(email string, body string) error {
	return nil
}

// RootDomain returns the MockMailer's rootdomain, always localhost
// TODO: this is probably broken
func (mm *MockMailer) RootDomain() string {
	return "http://localhost"
}

// StdoutMailer writes all emails to Stdout, perfect for dev / debugging
type StdoutMailer struct {
	Domain string
}

// Send writes the email to stdout
func (*StdoutMailer) Send(email string, body string) error {
	log.Println("hydrocarbon: new mail to", email, "\n", body)
	return nil
}

// RootDomain returns the StdoutMailer root domain
func (sm *StdoutMailer) RootDomain() string {
	return sm.Domain
}
