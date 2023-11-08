package rx

import (
	"context"
	"sync"
)

var waitGroupKey int

// A WaitGroup waits for a collection of goroutines to finish.
// The main goroutine calls Add to set the number of
// goroutines to wait for. Then each of the goroutines
// runs and calls Done when finished. At the same time,
// Wait can be used to block until all goroutines have finished.
//
// A WaitGroup must not be copied after first use.
type WaitGroup sync.WaitGroup

// WithWaitGroup returns a copy of ctx in which wg is associated.
func WithWaitGroup(ctx context.Context, wg *WaitGroup) context.Context {
	return context.WithValue(ctx, &waitGroupKey, wg)
}

// WaitGroupFromContext returns the WaitGroup that was associated in ctx.
//
// If there is no WaitGroup associated in ctx, WaitGroupFromContext returns
// nil; though, [WaitGroup.Go] is still safe for use.
func WaitGroupFromContext(ctx context.Context) *WaitGroup {
	wg, _ := ctx.Value(&waitGroupKey).(*WaitGroup)
	return wg
}

// Add adds delta, which may be negative, to the WaitGroup counter.
// If the counter becomes zero, all goroutines blocked on Wait are released.
// If the counter goes negative, Add panics.
func (wg *WaitGroup) Add(delta int) {
	(*sync.WaitGroup)(wg).Add(delta)
}

// Done decrements the WaitGroup counter by one.
func (wg *WaitGroup) Done() {
	(*sync.WaitGroup)(wg).Done()
}

// Wait blocks until the WaitGroup counter is zero.
func (wg *WaitGroup) Wait() {
	(*sync.WaitGroup)(wg).Wait()
}

// Go calls f in a new goroutine.
//
// If wg is not nil, Go calls wg.Add(1) before starting a goroutine
// to call f, and calls wg.Done() after f returns.
func (wg *WaitGroup) Go(f func()) {
	if wg == nil {
		go f()
		return
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		f()
	}()
}
