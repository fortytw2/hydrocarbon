package httpx

import (
	"io"
	"io/ioutil"
)

// DrainAndClose drains and closes an io.ReadCloser
func DrainAndClose(rc io.ReadCloser) error {
	_, err := io.Copy(ioutil.Discard, io.LimitReader(rc, 1024*8))
	if err != nil {
		return err
	}
	return rc.Close()
}
