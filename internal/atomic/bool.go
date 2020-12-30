package atomic

import (
	"sync/atomic"
)

// Bool is a type for atomic operations on bool.
type Bool struct {
	uint32
}

// FromBool creates an Bool with an initial value.
func FromBool(val bool) Bool {
	return Bool{boolToUint32(val)}
}

// Cas executes the compare-and-swap operation for *addr.
func (addr *Bool) Cas(old, new bool) bool {
	return atomic.CompareAndSwapUint32(&addr.uint32, boolToUint32(old), boolToUint32(new))
}

// Equals atomically loads *addr and checks if it equals to val.
func (addr *Bool) Equals(val bool) bool { return addr.Load() == val }

// False atomically loads *addr and checks if it's false.
func (addr *Bool) False() bool { return !addr.Load() }

// Load atomically loads *addr.
func (addr *Bool) Load() bool {
	return atomic.LoadUint32(&addr.uint32) != 0
}

// Store atomically stores val into *addr.
func (addr *Bool) Store(val bool) {
	atomic.StoreUint32(&addr.uint32, boolToUint32(val))
}

// Swap atomically stores val into *addr and returns the previous value.
func (addr *Bool) Swap(val bool) bool {
	return atomic.SwapUint32(&addr.uint32, boolToUint32(val)) != 0
}

// True atomically loads *addr and checks if it's true.
func (addr *Bool) True() bool { return addr.Load() }

func boolToUint32(val bool) uint32 {
	if val {
		return 1
	}

	return 0
}
