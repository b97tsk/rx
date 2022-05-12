package atomic

import (
	"sync/atomic"
)

// Bool is a type for atomic operations on bool.
type Bool struct {
	uint32
}

// FromBool creates an Bool with an initial value.
func FromBool(v bool) Bool {
	return Bool{boolToUint32(v)}
}

// CompareAndSwap executes the compare-and-swap operation for *addr.
func (addr *Bool) CompareAndSwap(old, new bool) bool {
	return atomic.CompareAndSwapUint32(&addr.uint32, boolToUint32(old), boolToUint32(new))
}

// Equal atomically loads *addr and checks if it equals to v.
func (addr *Bool) Equal(v bool) bool { return addr.Load() == v }

// False atomically loads *addr and checks if it's false.
func (addr *Bool) False() bool { return !addr.Load() }

// Load atomically loads *addr.
func (addr *Bool) Load() bool {
	return atomic.LoadUint32(&addr.uint32) != 0
}

// Store atomically stores v into *addr.
func (addr *Bool) Store(v bool) {
	atomic.StoreUint32(&addr.uint32, boolToUint32(v))
}

// Swap atomically stores v into *addr and returns the previous value.
func (addr *Bool) Swap(v bool) bool {
	return atomic.SwapUint32(&addr.uint32, boolToUint32(v)) != 0
}

// True atomically loads *addr and checks if it's true.
func (addr *Bool) True() bool { return addr.Load() }

func boolToUint32(v bool) uint32 {
	if v {
		return 1
	}

	return 0
}
