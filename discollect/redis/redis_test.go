//+build integration

package redis

import (
	"testing"

	"github.com/fortytw2/dockertest"
	"github.com/fortytw2/hydrocarbon/discollect"
)

func TestRedis(t *testing.T) {
	c, err := dockertest.RunContainer("redis:alpine", "6379", func(addr string) error {
		_, err := NewQueue(addr, 0)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Shutdown()

	q, err := NewQueue(c.Addr, 0)
	if err != nil {
		t.Fatal(q)
	}

	t.Run("standard", discollect.QueueTests(t, q, func() {
		err := q.ResetAll()
		if err != nil {
			panic(err)
		}
	}))
}
