// +build integration

package discollect

import (
	"testing"
)

func TestMemQueue(t *testing.T) {
	mq := NewMemQueue()

	t.Run("standard", QueueTests(t, mq, func() {
		mq.reset()
	}))
}
