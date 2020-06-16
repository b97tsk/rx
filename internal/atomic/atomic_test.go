package atomic_test

import (
	"testing"

	"github.com/b97tsk/rx/internal/atomic"
)

func TestUint32(t *testing.T) {
	u := atomic.Uint32(42)

	assert(t, u.Load() == 42, "Load didn't work.")
	assert(t, u.Add(8) == 50, "Add didn't work.")
	assert(t, u.Sub(5) == 45, "Sub didn't work.")

	assert(t, u.Cas(45, 54), "Cas didn't report a swap.")
	assert(t, u.Equals(54), "Cas didn't work.")

	assert(t, u.Swap(33) == 54, "Swap didn't return the old value.")
	assert(t, u.Equals(33), "Swap didn't work.")

	u.Store(42)
	assert(t, u.Equals(42), "Store didn't work.")
}

func assert(t *testing.T, ok bool, message string) {
	if !ok {
		t.Fatal(message)
	}
}
