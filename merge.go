package rx

import (
	"context"

	"github.com/b97tsk/rx/x/queue"
)

// A MergeConfigure is a configure for Merge.
type MergeConfigure struct {
	Project    func(interface{}, int) Observable
	Concurrent int
}

// Use creates an Operator from this configure.
func (configure MergeConfigure) Use() Operator {
	return func(source Observable) Observable {
		obs := mergeObservable{source, configure}
		return Create(obs.Subscribe)
	}
}

type mergeObservable struct {
	Source Observable
	MergeConfigure
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
	return FromObservables(observables...).Pipe(operators.MergeAll())
}

// MergeAll converts a higher-order Observable into a first-order Observable
// which concurrently delivers all values that are emitted on the inner
// Observables.
func (Operators) MergeAll() Operator {
	return operators.MergeMap(ProjectToObservable)
}

// MergeMap creates an Observable that projects each source value to an
// Observable which is merged in the output Observable.
//
// MergeMap maps each value to an Observable, then flattens all of these inner
// Observables using MergeAll.
func (Operators) MergeMap(project func(interface{}, int) Observable) Operator {
	return MergeConfigure{project, -1}.Use()
}

// MergeMapTo creates an Observable that projects each source value to the same
// Observable which is merged multiple times in the output Observable.
//
// It's like MergeMap, but maps each value always to the same inner Observable.
func (Operators) MergeMapTo(inner Observable) Operator {
	return operators.MergeMap(func(interface{}, int) Observable { return inner })
}
