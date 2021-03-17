package lockfree

import (
	"time"
	"unsafe"
)

// LKStack returns an empty queue.
func NewLKStack() *LKStack {
	n := unsafe.Pointer(&node{})
	return &LKStack{head: n}
}

// LKStack is a lock-free unbounded stack.
type LKStack struct {
	len  int64
	head unsafe.Pointer
}

func (q *LKStack) IsEmpty() bool {
	return q.Len() == 0
}

func (q *LKStack) Len() int64 {
	return q.len
}

// LKStack puts the given value v at the tail of the stack.
func (q *LKStack) Push(v interface{}) {
	n := &node{value: v}
	for {
		head := load(&q.head)
		next := load(&n.next)
		cas(&n.next, next, head)
		if cas(&q.head, head, n) {
			inrInt64(&q.len)
			return
		}

		time.Sleep(time.Nanosecond)
	}
}

// Pop removes and returns the value at the head of the stack.
// It returns nil if the stack is empty.
func (q *LKStack) Pop() interface{} {
	for {
		head := load(&q.head)
		next := load(&head.next)
		if next == nil { // is stack empty?
			return nil
		} else {
			// read value before CAS otherwise another Pop might free the next node
			v := head.value
			if cas(&q.head, head, next) {
				dcrInt64(&q.len)
				return v // Pop is done.  return
			}
		}
		time.Sleep(time.Nanosecond)
	}
}
