package queue

import (
	"testing"

	"github.com/fortytw2/discollect"
)

// TestQueue is a reusable conformance test that verifies the publically exposed interface
// of a given discollect.Queue
func TestQueue(t *testing.T, q discollect.Queue) {

	t.Run("basic-operations", basicQueueOps(q))
}

func basicQueueOps(q discollect.Queue) func(t *testing.T) {
	return func(t *testing.T) {

	}
}
