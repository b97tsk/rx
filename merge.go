package rx

import (
	"container/list"
	"context"
	"sync"
)

type mergeMapOperator struct {
	source     Operator
	project    func(interface{}, int) Observable
	concurrent int
}

func (op mergeMapOperator) ApplyOptions(options []Option) Operator {
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

func (op mergeMapOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()
	mu := sync.Mutex{}
	outerIndex := -1
	activeCount := 0
	buffer := list.List{}
	completeSignal := make(chan struct{}, 1)

	concurrent := op.concurrent
	if concurrent == 0 {
		concurrent = -1
	}

	var doNextLocked func()

	doNextLocked = func() {
		outerValue := buffer.Remove(buffer.Front())
		outerIndex++
		outerIndex := outerIndex

		// calls project synchronously
		obsv := op.project(outerValue, outerIndex)

		go obsv.Subscribe(ctx, ObserverFunc(func(t Notification) {
			switch {
			case t.HasValue:
				ob.Next(t.Value)

			case t.HasError:
				ob.Error(t.Value.(error))
				cancel()

			default:
				mu.Lock()

				if buffer.Len() > 0 {
					doNextLocked()
					mu.Unlock()
					break
				}

				activeCount--
				mu.Unlock()

				select {
				case completeSignal <- struct{}{}:
				default:
				}
			}
		}))
	}

	op.source.Call(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			mu.Lock()

			buffer.PushBack(t.Value)

			if activeCount != concurrent {
				activeCount++
				doNextLocked()
			}

			mu.Unlock()

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
					ob.Complete()
					cancel()
				}()
				return
			}
			mu.Unlock()
			ob.Complete()
			cancel()
		}
	}))

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
	op := mergeMapOperator{
		source:     o.Op,
		project:    projectToObservable,
		concurrent: -1,
	}
	return Observable{op}.Mutex()
}

// MergeMap creates an Observable that projects each source value to an
// Observable which is merged in the output Observable.
//
// MergeMap maps each value to an Observable, then flattens all of these inner
// Observables using MergeAll.
func (o Observable) MergeMap(project func(interface{}, int) Observable) Observable {
	op := mergeMapOperator{
		source:     o.Op,
		project:    project,
		concurrent: -1,
	}
	return Observable{op}.Mutex()
}

// MergeMapTo creates an Observable that projects each source value to the same
// Observable which is merged multiple times in the output Observable.
//
// It's like MergeMap, but maps each value always to the same inner Observable.
func (o Observable) MergeMapTo(inner Observable) Observable {
	op := mergeMapOperator{
		source:     o.Op,
		project:    func(interface{}, int) Observable { return inner },
		concurrent: -1,
	}
	return Observable{op}.Mutex()
}
