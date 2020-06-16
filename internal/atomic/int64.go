package atomic

import (
	"sync/atomic"
)

// Int64s is a type for atomic operations on int64.
type Int64s struct {
	int64
}

// Int64 creates an Int64s with an initial value.
func Int64(i int64) Int64s {
	return Int64s{i}
}

// Add atomically adds delta to *addr and returns the new value.
func (addr *Int64s) Add(delta int64) int64 {
	return atomic.AddInt64(&addr.int64, delta)
}

// Cas executes the compare-and-swap operation for *addr.
func (addr *Int64s) Cas(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&addr.int64, old, new)
}

// Equals atomically loads *addr and checks if it equals to val.
func (addr *Int64s) Equals(val int64) bool {
	return atomic.LoadInt64(&addr.int64) == val
}

// Load atomically loads *addr.
func (addr *Int64s) Load() int64 {
	return atomic.LoadInt64(&addr.int64)
}

// Store atomically stores val into *addr.
func (addr *Int64s) Store(val int64) {
	atomic.StoreInt64(&addr.int64, val)
}

// Sub atomically subtracts delta from *addr and returns the new value.
func (addr *Int64s) Sub(delta int64) int64 {
	return atomic.AddInt64(&addr.int64, ^(delta - 1))
}

// Swap atomically stores val into *addr and returns the previous value.
func (addr *Int64s) Swap(val int64) int64 {
	return atomic.SwapInt64(&addr.int64, val)
}
