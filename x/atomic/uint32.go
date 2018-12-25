package atomic

import (
	"sync/atomic"
)

// Uint32 is an uint32 type for atomic operations.
type Uint32 uint32

// Add atomically adds delta to *addr and returns the new value.
func (addr *Uint32) Add(delta uint32) uint32 {
	return atomic.AddUint32((*uint32)(addr), delta)
}

// Cas executes the compare-and-swap operation for *addr.
func (addr *Uint32) Cas(old, new uint32) bool {
	return atomic.CompareAndSwapUint32((*uint32)(addr), old, new)
}

// Equals atomically loads *addr and checks if it equals to val.
func (addr *Uint32) Equals(val uint32) bool {
	return atomic.LoadUint32((*uint32)(addr)) == val
}

// Load atomically loads *addr.
func (addr *Uint32) Load() uint32 {
	return atomic.LoadUint32((*uint32)(addr))
}

// Store atomically stores val into *addr.
func (addr *Uint32) Store(val uint32) {
	atomic.StoreUint32((*uint32)(addr), val)
}

// Sub atomically substracts delta from *addr and returns the new value.
func (addr *Uint32) Sub(delta uint32) uint32 {
	return atomic.AddUint32((*uint32)(addr), ^(delta - 1))
}

// Swap atomically stores val into *addr and returns the previous value.
func (addr *Uint32) Swap(val uint32) uint32 {
	return atomic.SwapUint32((*uint32)(addr), val)
}
