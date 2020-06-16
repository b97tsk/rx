package misc

import (
	"github.com/b97tsk/rx/internal/atomic"
)

// AvoidRecursion avoids to call some function recursively.
type AvoidRecursion struct {
	counter atomic.Uint32s
}

// Do calls a function in a non-recursive way when Do is called recursively.
func (ar *AvoidRecursion) Do(work func()) {
	if ar.counter.Add(1) > 1 {
		return
	}
	again := true
	for again {
		work()
		again = ar.counter.Sub(1) > 0
	}
}
