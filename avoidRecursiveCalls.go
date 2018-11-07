package rx

import (
	"math"
	"sync/atomic"
)

type avoidRecursiveCalls struct {
	n uint32
}

func (arc *avoidRecursiveCalls) Do(work func()) {
	if atomic.AddUint32(&arc.n, 1) > 1 {
		return
	}
	yes := true
	for yes {
		work()
		yes = atomic.AddUint32(&arc.n, math.MaxUint32) > 0
	}
}
