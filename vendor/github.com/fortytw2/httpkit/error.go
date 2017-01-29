package httpkit

// An Error is a error with a status
type Error interface {
	error
	Status() int
}

type httpErr struct {
	err    error
	status int
}

// Wrap adds a status code to an error
func Wrap(err error, status int) error {
	return &httpErr{
		err:    err,
		status: status,
	}
}

func (he httpErr) Error() string {
	return he.err.Error()
}

func (he httpErr) Status() int {
	return he.status
}
