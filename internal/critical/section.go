package critical

import (
	"sync"
)

type Section struct {
	mu     sync.Mutex
	closed bool
}

func Enter(s *Section) bool {
	s.mu.Lock()

	if s.closed {
		s.mu.Unlock()
		return false
	}

	return true
}

func Leave(s *Section) {
	s.mu.Unlock()
}

func Close(s *Section) {
	s.closed = true
	s.mu.Unlock()
}
