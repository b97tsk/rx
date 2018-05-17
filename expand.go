package rx

import (
	"container/list"
	"context"
	"sync"
)

type expandOperator struct {
	project    func(interface{}, int) Observable
	concurrent int
}

func (op expandOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
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

		ob.Next(outerValue)

		// calls project synchronously
		obsv := op.project(outerValue, outerIndex)

		go obsv.Subscribe(ctx, func(t Notification) {
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

// Expand recursively projects each source value to an Observable which is
// merged in the output Observable.
//
// It's similar to MergeMap, but applies the projection function to every
// source value as well as every output value. It's recursive.
func (o Observable) Expand(project func(interface{}, int) Observable) Observable {
	op := expandOperator{project, -1}
	return o.Lift(op.Call).Mutex()
}
