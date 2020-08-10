package norec

import (
	"sync/atomic"
)

// Wrap returns a function that calls f in a non-recursive way when it is
// called recursively.
func Wrap(f func()) func() {
	var n uint32
	return func() {
		if atomic.AddUint32(&n, 1) > 1 {
			return
		}
		again := true
		for again {
			f()
			again = atomic.AddUint32(&n, ^uint32(0)) > 0
		}
	}
}
