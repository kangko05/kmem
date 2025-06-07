package tests

import (
	"kmem/internal/queue"
	"testing"
)

func TestQueue(t *testing.T) {
	q := queue.New(t.Context())

	q.Add(queue.Q_TEST_ITEM)
}
