package rx

import (
	"context"

	"github.com/b97tsk/rx/x/queue"
)

type concatObservable struct {
	Source  Observable
	Project func(interface{}, int) Observable
}

func (obs concatObservable) Subscribe(ctx context.Context, sink Observer) {
	sink = Mutex(sink)

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
}

// Concat creates an output Observable which sequentially emits all values
// from given Observable and then moves on to the next.
//
// Concat concatenates multiple Observables together by sequentially emitting
// their values, one Observable after the other.
func Concat(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	obs := concatObservable{
		FromObservables(observables...),
		ProjectToObservable,
	}
	return Create(obs.Subscribe)
}
