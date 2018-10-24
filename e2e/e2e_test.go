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
	Run  func(*chromedp.Chrome, addr string, db *pg.DB, mm *hydrocarbon.MockMailer) error
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
			// TODO(fortytw2): does this work right?
			t.Parallel()

			// this will block if there are no free instances
			s, release := pool.getServer()
			defer release()

			chrome, chromeRelease := chromePool.getServer()
			defer chromeRelease()

			if c.Run != nil {
				err := c.Run(chrome, s.addr, s.db, s.mm)
				if err != nil {
					t.Fatal(err)
				}
			}

			s.db.TruncateTables()
		})
	}
}
