package rx

import (
	"context"
	"sync"

	"github.com/b97tsk/rx/x/queue"
)

type concatOperator struct {
	Project func(interface{}, int) Observable
}

func (op concatOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Mutex(Finally(sink, cancel))

	var (
		mutex        sync.Mutex
		outerIndex   = -1
		activeCount  = 1
		buffer       queue.Queue
		doNextLocked func()
	)

	doNextLocked = func() {
		var avoidRecursive avoidRecursiveCalls
		avoidRecursive.Do(func() {
			if buffer.Len() == 0 {
				activeCount--
				if activeCount == 0 {
					sink.Complete()
				}
				return
			}

			outerIndex++
			outerIndex := outerIndex
			outerValue := buffer.PopFront()

			obs := op.Project(outerValue, outerIndex)
			obs.Subscribe(ctx, func(t Notification) {
				switch {
				case t.HasValue || t.HasError:
					sink(t)
				default:
					avoidRecursive.Do(func() {
						mutex.Lock()
						doNextLocked()
						mutex.Unlock()
					})
				}
			})
		})
	}

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			mutex.Lock()
			buffer.PushBack(t.Value)
			if activeCount == 1 {
				activeCount++
				doNextLocked()
			}
			mutex.Unlock()

		case t.HasError:
			sink(t)

		default:
			mutex.Lock()
			activeCount--
			if activeCount == 0 {
				sink(t)
			}
			mutex.Unlock()
		}
	})

	return ctx, cancel
}

// Concat creates an output Observable which sequentially emits all values
// from given Observable and then moves on to the next.
//
// Concat concatenates multiple Observables together by sequentially emitting
// their values, one Observable after the other.
func Concat(observables ...Observable) Observable {
	return FromObservables(observables...).Pipe(operators.ConcatAll())
}

// ConcatAll converts a higher-order Observable into a first-order Observable
// by concatenating the inner Observables in order.
//
// ConcatAll flattens an Observable-of-Observables by putting one inner
// Observable after the other.
func (Operators) ConcatAll() OperatorFunc {
	return operators.ConcatMap(ProjectToObservable)
}

// ConcatMap projects each source value to an Observable which is merged in
// the output Observable, in a serialized fashion waiting for each one to
// complete before merging the next.
//
// ConcatMap maps each value to an Observable, then flattens all of these inner
// Observables using ConcatAll.
func (Operators) ConcatMap(project func(interface{}, int) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := concatOperator{project}
		return source.Lift(op.Call)
	}
}

// ConcatMapTo projects each source value to the same Observable which is
// merged multiple times in a serialized fashion on the output Observable.
//
// It's like ConcatMap, but maps each value always to the same inner Observable.
func (Operators) ConcatMapTo(inner Observable) OperatorFunc {
	return operators.ConcatMap(func(interface{}, int) Observable { return inner })
}
