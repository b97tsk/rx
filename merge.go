package rx

import (
	"context"
	"sync"

	"github.com/b97tsk/rx/x/queue"
)

// A MergeConfigure is a configure for Merge.
type MergeConfigure struct {
	Project    func(interface{}, int) Observable
	Concurrent int
}

// MakeFunc creates an OperatorFunc from this type.
func (conf MergeConfigure) MakeFunc() OperatorFunc {
	return MakeFunc(mergeOperator(conf).Call)
}

type mergeOperator MergeConfigure

func (op mergeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Mutex(Finally(sink, cancel))

	var (
		mutex           sync.Mutex
		outerIndex      = -1
		activeCount     int
		sourceCompleted bool
		buffer          queue.Queue
		doNextLocked    func()
	)

	concurrent := op.Concurrent
	if concurrent == 0 {
		concurrent = -1
	}

	doNextLocked = func() {
		outerIndex++
		outerIndex := outerIndex
		outerValue := buffer.PopFront()

		// calls op.Project synchronously
		obs := op.Project(outerValue, outerIndex)

		go obs.Subscribe(ctx, func(t Notification) {
			switch {
			case t.HasValue || t.HasError:
				sink(t)
			default:
				mutex.Lock()
				if buffer.Len() > 0 {
					doNextLocked()
				} else {
					activeCount--
					if activeCount == 0 && sourceCompleted {
						sink(t)
					}
				}
				mutex.Unlock()
			}
		})
	}

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			mutex.Lock()
			buffer.PushBack(t.Value)
			if activeCount != concurrent {
				activeCount++
				doNextLocked()
			}
			mutex.Unlock()

		case t.HasError:
			sink(t)

		default:
			mutex.Lock()
			sourceCompleted = true
			if activeCount == 0 {
				sink(t)
			}
			mutex.Unlock()
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
	return FromObservables(observables...).Pipe(operators.MergeAll())
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
		op := mergeOperator{project, -1}
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
