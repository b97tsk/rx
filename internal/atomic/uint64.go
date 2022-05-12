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

// CompareAndSwap executes the compare-and-swap operation for *addr.
func (addr *Uint64) CompareAndSwap(old, new uint64) bool {
	return atomic.CompareAndSwapUint64(&addr.uint64, old, new)
}

// Equal atomically loads *addr and checks if it equals to v.
func (addr *Uint64) Equal(v uint64) bool { return addr.Load() == v }

// Load atomically loads *addr.
func (addr *Uint64) Load() uint64 {
	return atomic.LoadUint64(&addr.uint64)
}

// Store atomically stores v into *addr.
func (addr *Uint64) Store(v uint64) {
	atomic.StoreUint64(&addr.uint64, v)
}

// Sub atomically subtracts delta from *addr and returns the new value.
func (addr *Uint64) Sub(delta uint64) uint64 {
	return atomic.AddUint64(&addr.uint64, ^(delta - 1))
}

// Swap atomically stores v into *addr and returns the previous value.
func (addr *Uint64) Swap(v uint64) uint64 {
	return atomic.SwapUint64(&addr.uint64, v)
}
