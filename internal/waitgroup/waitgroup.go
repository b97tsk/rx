package waitgroup

import (
	"context"
	"sync"
)

type WaitGroup = sync.WaitGroup

var ctxKey int

// Get returns the latest WaitGroup that was installed in ctx, or nil
// if there is none.
func Get(ctx context.Context) *WaitGroup {
	wg, _ := ctx.Value(&ctxKey).(*WaitGroup)
	return wg
}

// Hoist returns a copy of ctx that pulls the latest installed WaitGroup
// up close to the surface such that Get would perform better.
func Hoist(ctx context.Context) context.Context {
	return context.WithValue(ctx, &ctxKey, Get(ctx))
}

// Install returns a copy of ctx and a new WaitGroup that is bound to it.
//
// The WaitGroup returned has been added once, which means a Done call is
// expected somewhere after Install.
func Install(ctx context.Context) (context.Context, *WaitGroup) {
	var wg WaitGroup

	wg.Add(1)

	if pwg := Get(ctx); pwg != nil {
		pwg.Add(1)

		go func() {
			wg.Wait()
			pwg.Done()
		}()
	}

	return context.WithValue(ctx, &ctxKey, &wg), &wg
}
