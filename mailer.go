package kiasu

// Mailer is an interface for sending emails - used primarily for user confirmations
// but also send users notifications about "account expiry", etc
type Mailer interface{}

// FakeMailer returns a fake mailer suitable for using in tests / locally
func FakeMailer() Mailer {
	return nil
}
