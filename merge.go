package rx

import (
	"container/list"
	"context"
	"sync"
)

// MergeOperator is an operator type.
type MergeOperator struct {
	Project    func(interface{}, int) Observable
	Concurrent int
}

// MakeFunc creates an OperatorFunc from this operator.
func (op MergeOperator) MakeFunc() OperatorFunc {
	return MakeFunc(op.Call)
}

// Call invokes an execution of this operator.
func (op MergeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()
	sink = Mutex(sink)

	var (
		mu             sync.Mutex
		outerIndex     = -1
		activeCount    = 0
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
		outerIndex++
		outerIndex := outerIndex

		// calls op.Project synchronously
		obsv := op.Project(outerValue, outerIndex)

		go obsv.Subscribe(ctx, func(t Notification) {
			switch {
			case t.HasValue:
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
					sink.Complete()
					cancel()
				}()
				return
			}
			mu.Unlock()
			sink(t)
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
func Merge(observables ...Observable) Observable {
	return FromObservables(observables).Pipe(operators.MergeAll())
}

// MergeAll converts a higher-order Observable into a first-order Observable
// which concurrently delivers all values that are emitted on the inner
// Observables.
func (Operators) MergeAll() OperatorFunc {
	return operators.MergeMap(ProjectToObservable)
}

// MergeMap creates an Observable that projects each source value to an
// Observable which is merged in the output Observable.
//
// MergeMap maps each value to an Observable, then flattens all of these inner
// Observables using MergeAll.
func (Operators) MergeMap(project func(interface{}, int) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := MergeOperator{project, -1}
		return source.Lift(op.Call)
	}
}

// MergeMapTo creates an Observable that projects each source value to the same
// Observable which is merged multiple times in the output Observable.
//
// It's like MergeMap, but maps each value always to the same inner Observable.
func (Operators) MergeMapTo(inner Observable) OperatorFunc {
	return operators.MergeMap(func(interface{}, int) Observable { return inner })
}
