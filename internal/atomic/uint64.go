package atomic

import (
	"sync/atomic"
)

// Uint64 is a type for atomic operations on uint64.
type Uint64 struct {
	uint64
}

// FromUint64 creates an Uint64 with an initial value.
func FromUint64(u uint64) Uint64 {
	return Uint64{u}
}

// Add atomically adds delta to *addr and returns the new value.
func (addr *Uint64) Add(delta uint64) uint64 {
	return atomic.AddUint64(&addr.uint64, delta)
}

// Cas executes the compare-and-swap operation for *addr.
func (addr *Uint64) Cas(old, new uint64) bool {
	return atomic.CompareAndSwapUint64(&addr.uint64, old, new)
}

// Equals atomically loads *addr and checks if it equals to val.
func (addr *Uint64) Equals(val uint64) bool { return addr.Load() == val }

// Load atomically loads *addr.
func (addr *Uint64) Load() uint64 {
	return atomic.LoadUint64(&addr.uint64)
}

// Store atomically stores val into *addr.
func (addr *Uint64) Store(val uint64) {
	atomic.StoreUint64(&addr.uint64, val)
}

// Sub atomically subtracts delta from *addr and returns the new value.
func (addr *Uint64) Sub(delta uint64) uint64 {
	return atomic.AddUint64(&addr.uint64, ^(delta - 1))
}

// Swap atomically stores val into *addr and returns the previous value.
func (addr *Uint64) Swap(val uint64) uint64 {
	return atomic.SwapUint64(&addr.uint64, val)
}
