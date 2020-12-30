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

// Cas executes the compare-and-swap operation for *addr.
func (addr *Int64) Cas(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&addr.int64, old, new)
}

// Equals atomically loads *addr and checks if it equals to val.
func (addr *Int64) Equals(val int64) bool { return addr.Load() == val }

// Load atomically loads *addr.
func (addr *Int64) Load() int64 {
	return atomic.LoadInt64(&addr.int64)
}

// Store atomically stores val into *addr.
func (addr *Int64) Store(val int64) {
	atomic.StoreInt64(&addr.int64, val)
}

// Sub atomically subtracts delta from *addr and returns the new value.
func (addr *Int64) Sub(delta int64) int64 {
	return atomic.AddInt64(&addr.int64, -delta)
}

// Swap atomically stores val into *addr and returns the previous value.
func (addr *Int64) Swap(val int64) int64 {
	return atomic.SwapInt64(&addr.int64, val)
}
