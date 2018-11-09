package rx

import (
	"sync"
	"sync/atomic"
)

type cancellableLocker struct {
	mutex    sync.Mutex
	canceled uint32
}

func (l *cancellableLocker) Lock() bool {
	if atomic.LoadUint32(&l.canceled) != 0 {
		return false
	}

	l.mutex.Lock()

	if atomic.LoadUint32(&l.canceled) != 0 {
		l.mutex.Unlock()
		return false
	}

	return true
}

func (l *cancellableLocker) Unlock() {
	l.mutex.Unlock()
}

func (l *cancellableLocker) CancelAndUnlock() {
	atomic.StoreUint32(&l.canceled, 1)
	l.mutex.Unlock()
}
