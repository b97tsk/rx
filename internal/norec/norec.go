package norec

import (
	"sync/atomic"
)

// Wrap returns a function that calls f in a non-recursive way when it is
// called recursively.
func Wrap(f func()) func() {
	var n atomic.Uint32

	return func() {
		if n.Add(1) > 1 {
			return
		}

		for {
			f()

			if n.Add(^uint32(0)) == 0 {
				break
			}
		}
	}
}
