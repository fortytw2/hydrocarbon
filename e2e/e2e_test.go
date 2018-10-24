//+build integration

package e2e

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/pg"
)

var cases = []struct {
	Name    string
	Execute func(*chromedp.Chrome, addr string) error
	Verify  func(*chromedp.Chrome, addr string, db *pg.DB, mm *hydrocarbon.MockMailer) error
}{
	{
		// TODO: test cases go here
	},
}

func TestEndToEnd(t *testing.T) {

	instances := 1
	if val, ok := os.LookupEnv("HYDROCARBON_E2E_PARALLELISM"); ok {
		instances, err = strconv.Atoi(val)
		if err != nil {
			t.Fatal(err)
		}
	}

	pool := newTestServerPool(instances)
	chromePool := chromedp.NewPool(instances)

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			s, release := pool.getServer()
			defer release()

			chrome, chromeRelease := chromePool.getServer()
			defer chromeRelease()

			if c.Execute != nil {
				err := c.Execute(chrome, s.addr)
				if err != nil {
					t.Fatal(err)
				}
			}

			if c.Verify != nil {

				err := c.Verify(chrome, s.addr, s.db, s.mm)
				if err != nil {
					t.Fatal(err)
				}
			}

			s.db.TruncateTables()
		})
	}
}
