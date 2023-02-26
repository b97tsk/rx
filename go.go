package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/waitgroup"
)

// Go calls f in a new goroutine.
//
// Internally, a WaitGroup might have been installed in ctx.
// If that is the case, Go calls Add(1) on that WaitGroup before starting
// a goroutine to call f, and calls Done() on that WaitGroup after f returns.
//
// To correctly use Go, Go must be called by an Observable while being
// subscribed (by another Observable or Go), or by another Go.
func Go(ctx context.Context, f func()) {
	wg := waitgroup.Get(ctx)
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
