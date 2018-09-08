// +build integration

package discollect

import (
	"testing"
)

func TestMemQueue(t *testing.T) {
	mq := NewMemQueue()

	QueueTests(t, mq, func() {
		mq = NewMemQueue()
	})

}
