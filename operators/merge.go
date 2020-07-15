package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/queue"
)

// A MergeConfigure is a configure for Merge.
type MergeConfigure struct {
	Project    func(interface{}, int) rx.Observable
	Concurrent int
}

// Use creates an Operator from this configure.
func (configure MergeConfigure) Use() rx.Operator {
	if configure.Project == nil {
		configure.Project = projectToObservable
	}
	if configure.Concurrent == 0 {
		configure.Concurrent = -1
	}
	return func(source rx.Observable) rx.Observable {
		obs := mergeObservable{source, configure}
		return rx.Create(obs.Subscribe)
	}
}

type mergeObservable struct {
	Source rx.Observable
	MergeConfigure
}

func (obs mergeObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	sink = sink.Mutex()

	type X struct {
		Index           int
		Active          int
		Buffer          queue.Queue
		SourceCompleted bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var doNextLocked func(*X)

	doNextLocked = func(x *X) {
		sourceIndex := x.Index
		sourceValue := x.Buffer.Pop()
		x.Index++

		obs1 := obs.Project(sourceValue, sourceIndex)

		go obs1.Subscribe(ctx, func(t rx.Notification) {
			if t.HasValue || t.HasError {
				sink(t)
				return
			}
			x := <-cx
			if x.Buffer.Len() > 0 {
				doNextLocked(x)
			} else {
				x.Active--
				if x.Active == 0 && x.SourceCompleted {
					sink(t)
				}
			}
			cx <- x
		})
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			x := <-cx
			x.Buffer.Push(t.Value)
			if x.Active != obs.Concurrent {
				x.Active++
				doNextLocked(x)
			}
			cx <- x

		case t.HasError:
			sink(t)

		default:
			x := <-cx
			x.SourceCompleted = true
			if x.Active == 0 {
				sink(t)
			}
			cx <- x
		}
	})
}

// MergeAll converts a higher-order Observable into a first-order Observable
// which concurrently delivers all values that are emitted on the inner
// Observables.
func MergeAll() rx.Operator {
	return MergeMap(projectToObservable)
}

// MergeMap creates an Observable that projects each source value to an
// Observable which is merged in the output Observable.
//
// MergeMap maps each value to an Observable, then flattens all of these inner
// Observables using MergeAll.
func MergeMap(project func(interface{}, int) rx.Observable) rx.Operator {
	return MergeConfigure{project, -1}.Use()
}

// MergeMapTo creates an Observable that projects each source value to the same
// Observable which is merged multiple times in the output Observable.
//
// It's like MergeMap, but maps each value always to the same inner Observable.
func MergeMapTo(inner rx.Observable) rx.Operator {
	return MergeMap(func(interface{}, int) rx.Observable { return inner })
}
