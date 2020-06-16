package atomic

import (
	"sync/atomic"
)

// Uint64s is a type for atomic operations on uint64.
type Uint64s struct {
	uint64
}

// Uint64 creates an Uint64s with an initial value.
func Uint64(u uint64) Uint64s {
	return Uint64s{u}
}

// Add atomically adds delta to *addr and returns the new value.
func (addr *Uint64s) Add(delta uint64) uint64 {
	return atomic.AddUint64(&addr.uint64, delta)
}

// Cas executes the compare-and-swap operation for *addr.
func (addr *Uint64s) Cas(old, new uint64) bool {
	return atomic.CompareAndSwapUint64(&addr.uint64, old, new)
}

// Equals atomically loads *addr and checks if it equals to val.
func (addr *Uint64s) Equals(val uint64) bool {
	return atomic.LoadUint64(&addr.uint64) == val
}

// Load atomically loads *addr.
func (addr *Uint64s) Load() uint64 {
	return atomic.LoadUint64(&addr.uint64)
}

// Store atomically stores val into *addr.
func (addr *Uint64s) Store(val uint64) {
	atomic.StoreUint64(&addr.uint64, val)
}

// Sub atomically subtracts delta from *addr and returns the new value.
func (addr *Uint64s) Sub(delta uint64) uint64 {
	return atomic.AddUint64(&addr.uint64, ^(delta - 1))
}

// Swap atomically stores val into *addr and returns the previous value.
func (addr *Uint64s) Swap(val uint64) uint64 {
	return atomic.SwapUint64(&addr.uint64, val)
}
