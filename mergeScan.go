package rx

import (
	"container/list"
	"context"
	"sync"
)

type mergeScanOperator struct {
	source      Operator
	accumulator func(interface{}, interface{}) Observable
	seed        interface{}
	concurrent  int
}

func (op mergeScanOperator) ApplyOptions(options []Option) Operator {
	for _, opt := range options {
		switch t := opt.(type) {
		case concurrentOption:
			op.concurrent = t.Value
		default:
			panic(ErrUnsupportedOption)
		}
	}
	return op
}

func (op mergeScanOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()
	mu := sync.Mutex{}
	activeCount := 0
	seed := op.seed
	hasValue := false
	buffer := list.List{}
	completeSignal := make(chan struct{}, 1)

	concurrent := op.concurrent
	if concurrent == 0 {
		concurrent = -1
	}

	var doNextLocked func()

	doNextLocked = func() {
		outerValue := buffer.Remove(buffer.Front())

		// calls accumulator synchronously
		obsv := op.accumulator(seed, outerValue)

		go obsv.Subscribe(ctx, func(t Notification) {
			switch {
			case t.HasValue:
				mu.Lock()
				seed = t.Value
				hasValue = true
				mu.Unlock()

				ob.Next(t.Value)

			case t.HasError:
				ob.Error(t.Value.(error))
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

	op.source.Call(ctx, func(t Notification) {
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
			ob.Error(t.Value.(error))
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
						ob.Next(seed)
					}
					ob.Complete()
					cancel()
				}()
				return
			}
			mu.Unlock()
			if !hasValue {
				ob.Next(seed)
			}
			ob.Complete()
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
	op := mergeScanOperator{
		source:      o.Op,
		accumulator: accumulator,
		seed:        seed,
		concurrent:  -1,
	}
	return Observable{op}.Mutex()
}
