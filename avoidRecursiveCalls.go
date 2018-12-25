package rx

import (
	"github.com/b97tsk/rx/x/atomic"
)

type avoidRecursiveCalls struct {
	counter atomic.Uint32
}

func (arc *avoidRecursiveCalls) Do(work func()) {
	if arc.counter.Add(1) > 1 {
		return
	}
	yes := true
	for yes {
		work()
		yes = arc.counter.Sub(1) > 0
	}
}
