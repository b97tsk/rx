package atomic

import (
	"sync/atomic"
)

// Int32s is a type for atomic operations on int32.
type Int32s struct {
	int32
}

// Int32 creates an Int32s with an initial value.
func Int32(i int32) Int32s {
	return Int32s{i}
}

// Add atomically adds delta to *addr and returns the new value.
func (addr *Int32s) Add(delta int32) int32 {
	return atomic.AddInt32(&addr.int32, delta)
}

// Cas executes the compare-and-swap operation for *addr.
func (addr *Int32s) Cas(old, new int32) bool {
	return atomic.CompareAndSwapInt32(&addr.int32, old, new)
}

// Equals atomically loads *addr and checks if it equals to val.
func (addr *Int32s) Equals(val int32) bool {
	return atomic.LoadInt32(&addr.int32) == val
}

// Load atomically loads *addr.
func (addr *Int32s) Load() int32 {
	return atomic.LoadInt32(&addr.int32)
}

// Store atomically stores val into *addr.
func (addr *Int32s) Store(val int32) {
	atomic.StoreInt32(&addr.int32, val)
}

// Sub atomically subtracts delta from *addr and returns the new value.
func (addr *Int32s) Sub(delta int32) int32 {
	return atomic.AddInt32(&addr.int32, ^(delta - 1))
}

// Swap atomically stores val into *addr and returns the previous value.
func (addr *Int32s) Swap(val int32) int32 {
	return atomic.SwapInt32(&addr.int32, val)
}
