package lockfree_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/whiteCcinn/lockfree"
	"testing"
)

func TestStack(t *testing.T) {
	stack := lockfree.NewLKStack()
	loop := 10
	for i := 0; i < loop; i++ {
		stack.Push(i)
		assert.Equal(t, int64(i+1), stack.Len())
	}

	for i := 0; i < loop; i++ {
		assert.Equal(t, loop-i-1, stack.Pop())
	}
}
