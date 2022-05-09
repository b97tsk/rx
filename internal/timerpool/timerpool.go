package timerpool

import (
	"sync"
	"time"
)

var pool sync.Pool

// Get returns a time.Timer that will send the current time on its channel
// after at least duration d. Get might return an old time.Timer from pool.
func Get(d time.Duration) *time.Timer {
	if v := pool.Get(); v != nil {
		if t := v.(*time.Timer); !t.Reset(d) {
			return t
		}
	}

	return time.NewTimer(d)
}

// Put adds t to the pool. There must be no goroutines reading from t.C and
// whether t has expired must not be known.
func Put(t *time.Timer) {
	if t.Stop() {
		pool.Put(t)
	}
}

// Put adds t to the pool. There must be no goroutines reading from t.C and
// t must have expired.
func PutExpired(t *time.Timer) {
	pool.Put(t)
}
