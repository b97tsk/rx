package rx

import (
	"context"
	"sync"
)

// CongestingMergeOperator is an operator type.
type CongestingMergeOperator struct {
	Project    func(interface{}, int) Observable
	Concurrent int
}

// MakeFunc creates an OperatorFunc from this operator.
func (op CongestingMergeOperator) MakeFunc() OperatorFunc {
	return MakeFunc(op.Call)
}

// Call invokes an execution of this operator.
func (op CongestingMergeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	sink = Mutex(Finally(sink, cancel))

	var (
		mutex          sync.Mutex
		outerIndex     = -1
		activeCount    = 0
		completeSignal = make(chan struct{}, 1)
	)

	concurrent := op.Concurrent
	if concurrent == 0 {
		concurrent = -1
	}

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			mutex.Lock()
			for activeCount == concurrent {
				mutex.Unlock()
				select {
				case <-done:
					return
				case <-completeSignal:
				}
				mutex.Lock()
			}

			var try cancellableLocker
			defer func() {
				try.Lock()
				mutex.Unlock()
				try.CancelAndUnlock()
			}()

			activeCount++

			outerIndex++
			outerIndex := outerIndex
			outerValue := t.Value

			// calls op.Project synchronously
			obs := op.Project(outerValue, outerIndex)

			obs.Subscribe(ctx, func(t Notification) {
				switch {
				case t.HasValue || t.HasError:
					sink(t)
				default:
					if try.Lock() {
						activeCount--
						try.Unlock()
					} else {
						mutex.Lock()
						activeCount--
						mutex.Unlock()
					}
					select {
					case completeSignal <- struct{}{}:
					default:
					}
				}
			})

		case t.HasError:
			sink(t)

		default:
			mutex.Lock()
			if activeCount > 0 {
				go func() {
					for activeCount > 0 {
						mutex.Unlock()
						select {
						case <-done:
							return
						case <-completeSignal:
						}
						mutex.Lock()
					}
					mutex.Unlock()
					sink.Complete()
				}()
				return
			}
			mutex.Unlock()
			sink(t)
		}
	})

	return ctx, cancel
}

// CongestingMerge creates an output Observable which concurrently emits all
// values from every given input Observable.
//
// CongestingMerge flattens multiple Observables together by blending their
// values into one Observable.
//
// It's like Merge, but it may congest the source due to concurrent limit.
func CongestingMerge(observables ...Observable) Observable {
	return FromObservables(observables...).Pipe(operators.CongestingMergeAll())
}

// CongestingMergeAll converts a higher-order Observable into a first-order
// Observable which concurrently delivers all values that are emitted on the
// inner Observables.
//
// It's like MergeAll, but it may congest the source due to concurrent limit.
func (Operators) CongestingMergeAll() OperatorFunc {
	return operators.CongestingMergeMap(ProjectToObservable)
}

// CongestingMergeMap creates an Observable that projects each source value to
// an Observable which is merged in the output Observable.
//
// CongestingMergeMap maps each value to an Observable, then flattens all of
// these inner Observables using CongestingMergeAll.
//
// It's like MergeMap, but it may congest the source due to concurrent limit.
func (Operators) CongestingMergeMap(project func(interface{}, int) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := CongestingMergeOperator{project, -1}
		return source.Lift(op.Call)
	}
}

// CongestingMergeMapTo creates an Observable that projects each source value
// to the same Observable which is merged multiple times in the output
// Observable.
//
// It's like CongestingMergeMap, but maps each value always to the same inner
// Observable.
//
// It's like MergeMapTo, but it may congest the source due to concurrent limit.
func (Operators) CongestingMergeMapTo(inner Observable) OperatorFunc {
	return operators.CongestingMergeMap(func(interface{}, int) Observable { return inner })
}
