package atomic

import (
	"sync/atomic"
)

// Uint32 is a type for atomic operations on uint32.
type Uint32 struct {
	uint32
}

// FromUint32 creates an Uint32 with an initial value.
func FromUint32(u uint32) Uint32 {
	return Uint32{u}
}

// Add atomically adds delta to *addr and returns the new value.
func (addr *Uint32) Add(delta uint32) uint32 {
	return atomic.AddUint32(&addr.uint32, delta)
}

// CompareAndSwap executes the compare-and-swap operation for *addr.
func (addr *Uint32) CompareAndSwap(old, new uint32) bool {
	return atomic.CompareAndSwapUint32(&addr.uint32, old, new)
}

// Equal atomically loads *addr and checks if it equals to v.
func (addr *Uint32) Equal(v uint32) bool { return addr.Load() == v }

// Load atomically loads *addr.
func (addr *Uint32) Load() uint32 {
	return atomic.LoadUint32(&addr.uint32)
}

// Store atomically stores v into *addr.
func (addr *Uint32) Store(v uint32) {
	atomic.StoreUint32(&addr.uint32, v)
}

// Sub atomically subtracts delta from *addr and returns the new value.
func (addr *Uint32) Sub(delta uint32) uint32 {
	return atomic.AddUint32(&addr.uint32, ^(delta - 1))
}

// Swap atomically stores v into *addr and returns the previous value.
func (addr *Uint32) Swap(v uint32) uint32 {
	return atomic.SwapUint32(&addr.uint32, v)
}
