package rx

import (
	"context"

	"github.com/b97tsk/rx/x/queue"
)

type concatObservable struct {
	Source  Observable
	Project func(interface{}, int) Observable
}

func (obs concatObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Mutex(Finally(sink, cancel))

	type X struct {
		Index       int
		ActiveCount int
		Buffer      queue.Queue
	}
	cx := make(chan *X, 1)
	cx <- &X{ActiveCount: 1}

	var doNextLocked func(*X)

	doNextLocked = func(x *X) {
		var avoidRecursive avoidRecursiveCalls
		avoidRecursive.Do(func() {
			if x.Buffer.Len() == 0 {
				x.ActiveCount--
				if x.ActiveCount == 0 {
					sink.Complete()
				}
				return
			}

			sourceIndex := x.Index
			sourceValue := x.Buffer.PopFront()
			x.Index++

			obs := obs.Project(sourceValue, sourceIndex)
			obs.Subscribe(ctx, func(t Notification) {
				switch {
				case t.HasValue || t.HasError:
					sink(t)
				default:
					avoidRecursive.Do(func() {
						x := <-cx
						doNextLocked(x)
						cx <- x
					})
				}
			})
		})
	}

	obs.Source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			x := <-cx
			x.Buffer.PushBack(t.Value)
			if x.ActiveCount == 1 {
				x.ActiveCount++
				doNextLocked(x)
			}
			cx <- x

		case t.HasError:
			sink(t)

		default:
			x := <-cx
			x.ActiveCount--
			if x.ActiveCount == 0 {
				sink(t)
			}
			cx <- x
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
func (Operators) ConcatAll() Operator {
	return operators.ConcatMap(ProjectToObservable)
}

// ConcatMap projects each source value to an Observable which is merged in
// the output Observable, in a serialized fashion waiting for each one to
// complete before merging the next.
//
// ConcatMap maps each value to an Observable, then flattens all of these inner
// Observables using ConcatAll.
func (Operators) ConcatMap(project func(interface{}, int) Observable) Operator {
	return func(source Observable) Observable {
		return concatObservable{source, project}.Subscribe
	}
}

// ConcatMapTo projects each source value to the same Observable which is
// merged multiple times in a serialized fashion on the output Observable.
//
// It's like ConcatMap, but maps each value always to the same inner Observable.
func (Operators) ConcatMapTo(inner Observable) Operator {
	return operators.ConcatMap(func(interface{}, int) Observable { return inner })
}
