package httputil

import (
	"fmt"
	"io"
	"io/ioutil"
)

// DrainAndClose drains and closes an io.Closer
func DrainAndClose(rc io.ReadCloser) {
	_, err := io.Copy(ioutil.Discard, io.LimitReader(rc, 1024*8))
	if err != nil {
		fmt.Println("crit:", err)
	}
	err = rc.Close()
	if err != nil {
		fmt.Println("crit:", err)
	}
}
