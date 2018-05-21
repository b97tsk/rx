package rx

import (
	"container/list"
	"context"
	"sync"
)

type mergeScanOperator struct {
	Accumulator func(interface{}, interface{}) Observable
	Seed        interface{}
	Concurrent  int
}

func (op mergeScanOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	var (
		mu             sync.Mutex
		activeCount    = 0
		seed           = op.Seed
		hasValue       bool
		buffer         list.List
		completeSignal = make(chan struct{}, 1)
		doNextLocked   func()
	)

	concurrent := op.Concurrent
	if concurrent == 0 {
		concurrent = -1
	}

	doNextLocked = func() {
		outerValue := buffer.Remove(buffer.Front())

		// calls op.Accumulator synchronously
		obsv := op.Accumulator(seed, outerValue)

		go obsv.Subscribe(ctx, func(t Notification) {
			switch {
			case t.HasValue:
				mu.Lock()
				seed = t.Value
				hasValue = true
				mu.Unlock()

				sink(t)

			case t.HasError:
				sink(t)
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
			sink(t)
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
					if !hasValue {
						sink.Next(seed)
					}
					sink.Complete()
					cancel()
				}()
				return
			}
			mu.Unlock()
			if !hasValue {
				sink.Next(seed)
			}
			sink(t)
			cancel()
		}
	})

	return ctx, cancel
}

// MergeScan applies an accumulator function over the source Observable where
// the accumulator function itself returns an Observable, then each
// intermediate Observable returned is merged into the output Observable.
//
// It's like Scan, but the Observables returned by the accumulator are merged
// into the outer Observable.
func (o Observable) MergeScan(accumulator func(interface{}, interface{}) Observable, seed interface{}) Observable {
	op := mergeScanOperator{accumulator, seed, -1}
	return o.Lift(op.Call).Mutex()
}
