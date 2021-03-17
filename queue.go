package lockfree

import (
	"time"
	"unsafe"
)

// NewLKQueue returns an empty queue.
func NewLKQueue() *LKQueue {
	n := unsafe.Pointer(&node{})
	return &LKQueue{head: n, tail: n}
}

// LKQueue is a lock-free unbounded queue.
type LKQueue struct {
	len  int64
	head unsafe.Pointer
	tail unsafe.Pointer
}

type node struct {
	value interface{}
	next  unsafe.Pointer
}

// Enqueue puts the given value v at the tail of the queue.
func (q *LKQueue) Enqueue(v interface{}) {
	n := &node{value: v}
	for {
		tail := load(&q.tail)
		next := load(&tail.next)
		if tail == load(&q.tail) { // are tail and next consistent?
			if next == nil {
				if cas(&tail.next, next, n) {
					cas(&q.tail, tail, n) // Enqueue is done.  try to swing tail to the inserted node
					inrInt64(&q.len)
					return
				}
			} else { // tail was not pointing to the last node
				// try to swing Tail to the next node
				cas(&q.tail, tail, next)
			}
		}
		time.Sleep(time.Nanosecond)
	}
}

// Dequeue removes and returns the value at the head of the queue.
// It returns nil if the queue is empty.
func (q *LKQueue) Dequeue() interface{} {
	for {
		head := load(&q.head)
		tail := load(&q.tail)
		next := load(&head.next)
		if head == tail { // is queue empty or tail falling behind?
			if next == nil { // is queue empty?
				return nil
			}
			// tail is falling behind.  try to advance it
			cas(&q.tail, tail, next)
		} else {
			// read value before CAS otherwise another dequeue might free the next node
			v := next.value
			if cas(&q.head, head, next) {
				dcrInt64(&q.len)
				return v // Dequeue is done.  return
			}
		}
		time.Sleep(time.Nanosecond)
	}
}

func (q *LKQueue) Len() int64 {
	return q.len
}

func (q *LKQueue) IsEmpty() bool {
	return q.Len() == 0
}
