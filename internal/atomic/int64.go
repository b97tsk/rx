package atomic

import (
	"sync/atomic"
)

// Int64 is a type for atomic operations on int64.
type Int64 struct {
	int64
}

// FromInt64 creates an Int64 with an initial value.
func FromInt64(i int64) Int64 {
	return Int64{i}
}

// Add atomically adds delta to *addr and returns the new value.
func (addr *Int64) Add(delta int64) int64 {
	return atomic.AddInt64(&addr.int64, delta)
}

// CompareAndSwap executes the compare-and-swap operation for *addr.
func (addr *Int64) CompareAndSwap(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&addr.int64, old, new)
}

// Equal atomically loads *addr and checks if it equals to v.
func (addr *Int64) Equal(v int64) bool { return addr.Load() == v }

// Load atomically loads *addr.
func (addr *Int64) Load() int64 {
	return atomic.LoadInt64(&addr.int64)
}

// Store atomically stores v into *addr.
func (addr *Int64) Store(v int64) {
	atomic.StoreInt64(&addr.int64, v)
}

// Sub atomically subtracts delta from *addr and returns the new value.
func (addr *Int64) Sub(delta int64) int64 {
	return atomic.AddInt64(&addr.int64, -delta)
}

// Swap atomically stores v into *addr and returns the previous value.
func (addr *Int64) Swap(v int64) int64 {
	return atomic.SwapInt64(&addr.int64, v)
}
