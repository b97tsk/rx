package atomic

import (
	"sync/atomic"
)

// Uint32s is a type for atomic operations on uint32.
type Uint32s struct {
	uint32
}

// Uint32 creates an Uint32s with an initial value.
func Uint32(u uint32) Uint32s {
	return Uint32s{u}
}

// Add atomically adds delta to *addr and returns the new value.
func (addr *Uint32s) Add(delta uint32) uint32 {
	return atomic.AddUint32(&addr.uint32, delta)
}

// Cas executes the compare-and-swap operation for *addr.
func (addr *Uint32s) Cas(old, new uint32) bool {
	return atomic.CompareAndSwapUint32(&addr.uint32, old, new)
}

// Equals atomically loads *addr and checks if it equals to val.
func (addr *Uint32s) Equals(val uint32) bool {
	return atomic.LoadUint32(&addr.uint32) == val
}

// Load atomically loads *addr.
func (addr *Uint32s) Load() uint32 {
	return atomic.LoadUint32(&addr.uint32)
}

// Store atomically stores val into *addr.
func (addr *Uint32s) Store(val uint32) {
	atomic.StoreUint32(&addr.uint32, val)
}

// Sub atomically subtracts delta from *addr and returns the new value.
func (addr *Uint32s) Sub(delta uint32) uint32 {
	return atomic.AddUint32(&addr.uint32, ^(delta - 1))
}

// Swap atomically stores val into *addr and returns the previous value.
func (addr *Uint32s) Swap(val uint32) uint32 {
	return atomic.SwapUint32(&addr.uint32, val)
}
