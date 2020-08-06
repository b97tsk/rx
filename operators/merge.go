package operators

import (
	"context"
	"sync"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/queue"
)

// A MergeConfigure is a configure for Merge.
type MergeConfigure struct {
	Project    func(interface{}, int) rx.Observable
	Concurrent int
}

// Make creates an Operator from this configure.
func (configure MergeConfigure) Make() rx.Operator {
	if configure.Project == nil {
		configure.Project = projectToObservable
	}
	if configure.Concurrent == 0 {
		configure.Concurrent = -1
	}
	return func(source rx.Observable) rx.Observable {
		return mergeObservable{source, configure}.Subscribe
	}
}

type mergeObservable struct {
	Source rx.Observable
	MergeConfigure
}

func (obs mergeObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.WithCancel(cancel).Mutex()

	x := struct {
		sync.Mutex
		Queue     queue.Queue
		Index     int
		Workers   int
		Completed bool
	}{Index: -1}

	var subscribeLocked func()
	subscribeLocked = func() {
		x.Index++
		sourceIndex := x.Index
		sourceValue := x.Queue.Pop()
		obs1 := obs.Project(sourceValue, sourceIndex)
		go obs1.Subscribe(ctx, func(t rx.Notification) {
			if t.HasValue || t.HasError {
				sink(t)
				return
			}
			x.Lock()
			defer x.Unlock()
			if x.Queue.Len() > 0 {
				subscribeLocked()
			} else {
				x.Workers--
				if x.Completed && x.Workers == 0 {
					sink(t)
				}
			}
		})
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			x.Lock()
			defer x.Unlock()
			x.Queue.Push(t.Value)
			if x.Workers != obs.Concurrent {
				x.Workers++
				subscribeLocked()
			}

		case t.HasError:
			sink(t)

		default:
			x.Lock()
			defer x.Unlock()
			x.Completed = true
			if x.Workers == 0 {
				sink(t)
			}
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
	return MergeConfigure{project, -1}.Make()
}

// MergeMapTo creates an Observable that projects each source value to the same
// Observable which is merged multiple times in the output Observable.
//
// It's like MergeMap, but maps each value always to the same inner Observable.
func MergeMapTo(inner rx.Observable) rx.Operator {
	return MergeMap(func(interface{}, int) rx.Observable { return inner })
}
