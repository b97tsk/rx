package critical

import (
	"sync"
)

// A Section represents a critical section that only one goroutine can
// execute codes between a successful call of [Enter] and a call of [Leave]
// or [Close], at the same time.
type Section struct {
	mu     sync.Mutex
	closed bool
}

// Enter locks s and returns true.
// If s was closed, Enter returns false, and s remains unlocked.
func Enter(s *Section) bool {
	s.mu.Lock()

	if s.closed {
		s.mu.Unlock()
		return false
	}

	return true
}

// Leave unlocks s.
func Leave(s *Section) {
	s.mu.Unlock()
}

// Close closes and unlocks s.
// Future Enter calls return false after Close.
func Close(s *Section) {
	s.closed = true
	s.mu.Unlock()
}
