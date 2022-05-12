package atomic_test

import (
	"testing"

	"github.com/b97tsk/rx/internal/atomic"
)

func TestBool(t *testing.T) {
	t.Parallel()

	b := atomic.FromBool(true)

	assert(t, b.Load(), "Load didn't work.")

	assert(t, b.CompareAndSwap(true, false), "CompareAndSwap didn't report a swap.")
	assert(t, b.False(), "CompareAndSwap didn't work.")

	assert(t, !b.Swap(true), "Swap didn't return the old value.")
	assert(t, b.True(), "Swap didn't work.")

	b.Store(false)
	assert(t, b.Equal(false), "Store didn't work.")
}

func TestInt32(t *testing.T) {
	t.Parallel()

	i := atomic.FromInt32(42)

	assert(t, i.Load() == 42, "Load didn't work.")
	assert(t, i.Add(8) == 50, "Add didn't work.")
	assert(t, i.Sub(5) == 45, "Sub didn't work.")

	assert(t, i.CompareAndSwap(45, 54), "CompareAndSwap didn't report a swap.")
	assert(t, i.Equal(54), "CompareAndSwap didn't work.")

	assert(t, i.Swap(33) == 54, "Swap didn't return the old value.")
	assert(t, i.Equal(33), "Swap didn't work.")

	i.Store(42)
	assert(t, i.Equal(42), "Store didn't work.")
}

func TestInt64(t *testing.T) {
	t.Parallel()

	i := atomic.FromInt64(42)

	assert(t, i.Load() == 42, "Load didn't work.")
	assert(t, i.Add(8) == 50, "Add didn't work.")
	assert(t, i.Sub(5) == 45, "Sub didn't work.")

	assert(t, i.CompareAndSwap(45, 54), "CompareAndSwap didn't report a swap.")
	assert(t, i.Equal(54), "CompareAndSwap didn't work.")

	assert(t, i.Swap(33) == 54, "Swap didn't return the old value.")
	assert(t, i.Equal(33), "Swap didn't work.")

	i.Store(42)
	assert(t, i.Equal(42), "Store didn't work.")
}

func TestUint32(t *testing.T) {
	t.Parallel()

	u := atomic.FromUint32(42)

	assert(t, u.Load() == 42, "Load didn't work.")
	assert(t, u.Add(8) == 50, "Add didn't work.")
	assert(t, u.Sub(5) == 45, "Sub didn't work.")

	assert(t, u.CompareAndSwap(45, 54), "CompareAndSwap didn't report a swap.")
	assert(t, u.Equal(54), "CompareAndSwap didn't work.")

	assert(t, u.Swap(33) == 54, "Swap didn't return the old value.")
	assert(t, u.Equal(33), "Swap didn't work.")

	u.Store(42)
	assert(t, u.Equal(42), "Store didn't work.")
}

func TestUint64(t *testing.T) {
	t.Parallel()

	u := atomic.FromUint64(42)

	assert(t, u.Load() == 42, "Load didn't work.")
	assert(t, u.Add(8) == 50, "Add didn't work.")
	assert(t, u.Sub(5) == 45, "Sub didn't work.")

	assert(t, u.CompareAndSwap(45, 54), "CompareAndSwap didn't report a swap.")
	assert(t, u.Equal(54), "CompareAndSwap didn't work.")

	assert(t, u.Swap(33) == 54, "Swap didn't return the old value.")
	assert(t, u.Equal(33), "Swap didn't work.")

	u.Store(42)
	assert(t, u.Equal(42), "Store didn't work.")
}

func assert(t *testing.T, ok bool, message string) {
	t.Helper()

	if ok {
		return
	}

	t.Fatal(message)
}
