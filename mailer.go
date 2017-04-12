package hydrocarbon

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
