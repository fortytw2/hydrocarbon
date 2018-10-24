//+build integration

package e2e

import (
	"os"
	"runtime"
	"strconv"
	"testing"

	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/pg"
)

type E2ETestCase struct {
	Name string
	Run  func(browser *chromedp.Chrome, addr string, db *pg.DB, mm *hydrocarbon.MockMailer) error
}

var cases = []E2ETestCase{
	{
		// TODO: test cases go here
	},
}

func TestEndToEnd(t *testing.T) {
	var instances int
	if val, ok := os.LookupEnv("E2E_PARALLELISM"); ok {
		instances, err = strconv.Atoi(val)
		if err != nil {
			t.Fatal(err)
		}
	} else {
		// this should be a sensible default
		instances = runtime.NumCPU() / 2
	}

	pool := newTestServerPool(instances)
	chromePool := chromedp.NewPool(instances)

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			// TODO(fortytw2): does this work right with subtests
			t.Parallel()

			// this will block if there are no free instances
			s, release := pool.getServer()
			defer release()

			chrome, chromeRelease := chromePool.GetInstance()
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
