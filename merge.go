package rx

import (
	"context"

	"github.com/b97tsk/rx/x/queue"
)

type mergeObservable struct {
	Source     Observable
	Project    func(interface{}, int) Observable
	Concurrent int
}

func (obs mergeObservable) Subscribe(ctx context.Context, sink Observer) {
	sink = Mutex(sink)

	type X struct {
		Index           int
		ActiveCount     int
		SourceCompleted bool
		Buffer          queue.Queue
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var doNextLocked func(*X)

	doNextLocked = func(x *X) {
		sourceIndex := x.Index
		sourceValue := x.Buffer.PopFront()
		x.Index++

		// calls obs.Project synchronously
		obs := obs.Project(sourceValue, sourceIndex)

		go obs.Subscribe(ctx, func(t Notification) {
			switch {
			case t.HasValue || t.HasError:
				sink(t)
			default:
				x := <-cx
				if x.Buffer.Len() > 0 {
					doNextLocked(x)
				} else {
					x.ActiveCount--
					if x.ActiveCount == 0 && x.SourceCompleted {
						sink(t)
					}
				}
				cx <- x
			}
		})
	}

	obs.Source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			x := <-cx
			x.Buffer.PushBack(t.Value)
			if x.ActiveCount != obs.Concurrent {
				x.ActiveCount++
				doNextLocked(x)
			}
			cx <- x

		case t.HasError:
			sink(t)

		default:
			x := <-cx
			x.SourceCompleted = true
			if x.ActiveCount == 0 {
				sink(t)
			}
			cx <- x
		}
	})
}

// Merge creates an output Observable which concurrently emits all values from
// every given input Observable.
//
// Merge flattens multiple Observables together by blending their values into
// one Observable.
func Merge(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	obs := mergeObservable{
		FromObservables(observables...),
		ProjectToObservable,
		-1,
	}
	return Create(obs.Subscribe)
}
