package rx

import (
	"container/list"
	"context"
	"sync"
)

type mergeMapOperator struct {
	project    func(interface{}, int) Observable
	concurrent int
}

func (op mergeMapOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	var (
		mu             sync.Mutex
		outerIndex     = -1
		activeCount    = 0
		buffer         list.List
		completeSignal = make(chan struct{}, 1)
		doNextLocked   func()
	)

	concurrent := op.concurrent
	if concurrent == 0 {
		concurrent = -1
	}

	doNextLocked = func() {
		outerValue := buffer.Remove(buffer.Front())
		outerIndex++
		outerIndex := outerIndex

		// calls project synchronously
		obsv := op.project(outerValue, outerIndex)

		go obsv.Subscribe(ctx, func(t Notification) {
			switch {
			case t.HasValue:
				t.Observe(ob)

			case t.HasError:
				t.Observe(ob)
				cancel()

			default:
				mu.Lock()

				if buffer.Len() > 0 {
					defer mu.Unlock()
					doNextLocked()
					break
				}

				activeCount--
				mu.Unlock()

				select {
				case completeSignal <- struct{}{}:
				default:
				}
			}
		})
	}

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			mu.Lock()
			defer mu.Unlock()

			buffer.PushBack(t.Value)

			if activeCount != concurrent {
				activeCount++
				doNextLocked()
			}

		case t.HasError:
			t.Observe(ob)
			cancel()

		default:
			mu.Lock()
			if activeCount > 0 {
				go func() {
					for activeCount > 0 {
						mu.Unlock()
						select {
						case <-done:
							return
						case <-completeSignal:
						}
						mu.Lock()
					}
					mu.Unlock()
					ob.Complete()
					cancel()
				}()
				return
			}
			mu.Unlock()
			ob.Complete()
			cancel()
		}
	})

	return ctx, cancel
}

// Merge creates an output Observable which concurrently emits all values from
// every given input Observable.
//
// Merge flattens multiple Observables together by blending their values into
// one Observable.
func Merge(observables ...interface{}) Observable {
	return FromSlice(observables).MergeAll()
}

// MergeAll converts a higher-order Observable into a first-order Observable
// which concurrently delivers all values that are emitted on the inner
// Observables.
func (o Observable) MergeAll() Observable {
	op := mergeMapOperator{projectToObservable, -1}
	return o.Lift(op.Call).Mutex()
}

// MergeMap creates an Observable that projects each source value to an
// Observable which is merged in the output Observable.
//
// MergeMap maps each value to an Observable, then flattens all of these inner
// Observables using MergeAll.
func (o Observable) MergeMap(project func(interface{}, int) Observable) Observable {
	op := mergeMapOperator{project, -1}
	return o.Lift(op.Call).Mutex()
}

// MergeMapTo creates an Observable that projects each source value to the same
// Observable which is merged multiple times in the output Observable.
//
// It's like MergeMap, but maps each value always to the same inner Observable.
func (o Observable) MergeMapTo(inner Observable) Observable {
	return o.MergeMap(func(interface{}, int) Observable { return inner })
}
