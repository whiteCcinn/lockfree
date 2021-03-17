package lockfree_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/whiteCcinn/lockfree"
	"testing"
)

func TestQueue(t *testing.T) {
	queue := lockfree.NewLKQueue()
	loop := 10
	for i := 0; i < loop; i++ {
		queue.Enqueue(i)
		assert.Equal(t, int64(i+1), queue.Len())
	}

	for i := 0; i < loop; i++ {
		assert.Equal(t, i, queue.Dequeue())
	}
}
