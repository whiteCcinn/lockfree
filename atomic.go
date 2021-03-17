package lockfree

import (
	"sync/atomic"
	"time"
	"unsafe"
)

// load from atomic load pointer node
func load(p *unsafe.Pointer) (n *node) {
	return (*node)(atomic.LoadPointer(p))
}

// cas swap set
func cas(p *unsafe.Pointer, old, new *node) (ok bool) {
	return atomic.CompareAndSwapPointer(
		p, unsafe.Pointer(old), unsafe.Pointer(new))
}

// inrInt64 Increase
func inrInt64(i *int64) {
	t := int64(+1)
	for {
		value := atomic.LoadInt64(i)
		if atomic.CompareAndSwapInt64(i, value, value+t) {
			return
		}
		time.Sleep(time.Nanosecond)
	}
}

// dcrInt64 Decrease
func dcrInt64(i *int64) {
	t := int64(-1)
	for {
		value := atomic.LoadInt64(i)
		if atomic.CompareAndSwapInt64(i, value, value+t) {
			return
		}
		time.Sleep(time.Nanosecond)
	}
}
