package hydrocarbon

// A Mailer sends mail
type Mailer interface {
	Send(email string, body string) error
}
