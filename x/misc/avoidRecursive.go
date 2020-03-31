package misc

import (
	"github.com/b97tsk/rx/x/atomic"
)

// AvoidRecursive avoids to call some function recursively.
type AvoidRecursive struct {
	counter atomic.Uint32
}

// Do calls a specific function in a non-recursive way when it is called
// recursively.
func (ar *AvoidRecursive) Do(work func()) {
	if ar.counter.Add(1) > 1 {
		return
	}
	yes := true
	for yes {
		work()
		yes = ar.counter.Sub(1) > 0
	}
}
