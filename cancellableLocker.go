package rx

import (
	"sync"

	"github.com/b97tsk/rx/x/atomic"
)

type cancellableLocker struct {
	mutex    sync.Mutex
	canceled atomic.Uint32
}

func (l *cancellableLocker) Lock() bool {
	if l.canceled.Equals(1) {
		return false
	}

	l.mutex.Lock()

	if l.canceled.Equals(1) {
		l.mutex.Unlock()
		return false
	}

	return true
}

func (l *cancellableLocker) Unlock() {
	l.mutex.Unlock()
}

func (l *cancellableLocker) CancelAndUnlock() {
	l.canceled.Store(1)
	l.mutex.Unlock()
}
