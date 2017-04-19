package hydrocarbon

import "log"

// A Mailer sends mail
type Mailer interface {
	Send(email string, body string) error
}

type MockMailer struct {
	Mails []string
}

func (mm *MockMailer) Send(email string, body string) error {
	return nil
}

type StdoutMailer struct{}

func (*StdoutMailer) Send(email string, body string) error {
	log.Println("hydrocarbon: new mail to", email, "\n", body)
	return nil
}
