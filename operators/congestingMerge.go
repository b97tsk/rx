package operators

import (
	"context"
	"sync"

	"github.com/b97tsk/rx"
)

// A CongestingMergeConfigure is a configure for CongestingMerge.
type CongestingMergeConfigure struct {
	Project    func(interface{}, int) rx.Observable
	Concurrent int
}

// Make creates an Operator from this configure.
func (configure CongestingMergeConfigure) Make() rx.Operator {
	if configure.Project == nil {
		configure.Project = projectToObservable
	}
	if configure.Concurrent == 0 {
		configure.Concurrent = -1
	}
	return func(source rx.Observable) rx.Observable {
		return congestingMergeObservable{source, configure}.Subscribe
	}
}

type congestingMergeObservable struct {
	Source rx.Observable
	CongestingMergeConfigure
}

func (obs congestingMergeObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.WithCancel(cancel).Mutex()

	done := ctx.Done()
	complete := make(chan struct{}, 1)

	x := struct {
		sync.Mutex
		Index     int
		Workers   int
		Completed bool
	}{Index: -1}

	var observer rx.Observer

	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			x.Lock()
			for x.Workers == obs.Concurrent {
				x.Unlock()
				select {
				case <-done:
					observer = rx.Noop
					return
				case <-complete:
				}
				x.Lock()
			}
			x.Index++
			x.Workers++
			x.Unlock()

			obs1 := obs.Project(t.Value, x.Index)
			go obs1.Subscribe(ctx, func(t rx.Notification) {
				if t.HasValue || t.HasError {
					sink(t)
					return
				}
				x.Lock()
				defer x.Unlock()
				x.Workers--
				if x.Completed && x.Workers == 0 {
					sink(t)
				}
				select {
				case complete <- struct{}{}:
				default:
				}
			})

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
	}

	obs.Source.Subscribe(ctx, observer.Sink)
}

// CongestingMergeAll converts a higher-order Observable into a first-order
// Observable which concurrently delivers all values that are emitted on the
// inner Observables.
//
// It's like MergeAll, but it may congest the source due to concurrent limit.
func CongestingMergeAll() rx.Operator {
	return CongestingMergeMap(projectToObservable)
}

// CongestingMergeMap creates an Observable that projects each source value to
// an Observable which is merged in the output Observable.
//
// CongestingMergeMap maps each value to an Observable, then flattens all of
// these inner Observables using CongestingMergeAll.
//
// It's like MergeMap, but it may congest the source due to concurrent limit.
func CongestingMergeMap(project func(interface{}, int) rx.Observable) rx.Operator {
	return CongestingMergeConfigure{project, -1}.Make()
}

// CongestingMergeMapTo creates an Observable that projects each source value
// to the same Observable which is merged multiple times in the output
// Observable.
//
// It's like CongestingMergeMap, but maps each value always to the same inner
// Observable.
//
// It's like MergeMapTo, but it may congest the source due to concurrent limit.
func CongestingMergeMapTo(inner rx.Observable) rx.Operator {
	return CongestingMergeMap(func(interface{}, int) rx.Observable { return inner })
}
