package atomic

import (
	"sync/atomic"
)

// Int32 is a type for atomic operations on int32.
type Int32 struct {
	int32
}

// FromInt32 creates an Int32 with an initial value.
func FromInt32(i int32) Int32 {
	return Int32{i}
}

// Add atomically adds delta to *addr and returns the new value.
func (addr *Int32) Add(delta int32) int32 {
	return atomic.AddInt32(&addr.int32, delta)
}

// CompareAndSwap executes the compare-and-swap operation for *addr.
func (addr *Int32) CompareAndSwap(old, new int32) bool {
	return atomic.CompareAndSwapInt32(&addr.int32, old, new)
}

// Equal atomically loads *addr and checks if it equals to v.
func (addr *Int32) Equal(v int32) bool { return addr.Load() == v }

// Load atomically loads *addr.
func (addr *Int32) Load() int32 {
	return atomic.LoadInt32(&addr.int32)
}

// Store atomically stores v into *addr.
func (addr *Int32) Store(v int32) {
	atomic.StoreInt32(&addr.int32, v)
}

// Sub atomically subtracts delta from *addr and returns the new value.
func (addr *Int32) Sub(delta int32) int32 {
	return atomic.AddInt32(&addr.int32, -delta)
}

// Swap atomically stores v into *addr and returns the previous value.
func (addr *Int32) Swap(v int32) int32 {
	return atomic.SwapInt32(&addr.int32, v)
}
