package rx

import (
	"sync"
	"sync/atomic"
)

type cancellableLocker struct {
	mu       sync.Mutex
	canceled uint32
}

func (l *cancellableLocker) Lock() bool {
	if atomic.LoadUint32(&l.canceled) != 0 {
		return false
	}

	l.mu.Lock()

	if atomic.LoadUint32(&l.canceled) != 0 {
		l.mu.Unlock()
		return false
	}

	return true
}

func (l *cancellableLocker) Unlock() {
	l.mu.Unlock()
}

func (l *cancellableLocker) CancelAndUnlock() {
	atomic.StoreUint32(&l.canceled, 1)
	l.mu.Unlock()
}
