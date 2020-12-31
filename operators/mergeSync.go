package operators

import (
	"context"
	"sync"

	"github.com/b97tsk/rx"
)

// MergeSyncAll converts a higher-order Observable into a first-order
// Observable which concurrently delivers all values that are emitted on
// the inner Observables.
//
// It's like MergeAll, but it does not buffer the source.
//
// Concurrency should be positive, otherwise it works just like MergeAll.
func MergeSyncAll(concurrency int) rx.Operator {
	return MergeSyncMap(projectToObservable, concurrency)
}

// MergeSyncMap projects each source value to an Observable which is merged
// in the output Observable.
//
// MergeSyncMap maps each value to an Observable, then flattens all of these
// inner Observables using MergeSyncAll.
//
// It's like MergeMap, but it does not buffer the source.
//
// Concurrency should be positive, otherwise it works just like MergeMap.
func MergeSyncMap(
	project func(interface{}, int) rx.Observable,
	concurrency int,
) rx.Operator {
	if project == nil {
		panic("MergeSyncMap: project is nil")
	}

	if concurrency == 0 {
		panic("MergeSyncMap: concurrency is zero")
	}

	return func(source rx.Observable) rx.Observable {
		return mergeSyncObservable{source, project, concurrency}.Subscribe
	}
}

// MergeSyncMapTo projects each source value to the same Observable which
// is merged multiple times in the output Observable.
//
// It's like MergeMapSync, but maps each value always to the same inner
// Observable.
//
// It's like MergeMapTo, but it does not buffer the source.
//
// Concurrency should be positive, otherwise it works just like MergeMapTo.
func MergeSyncMapTo(inner rx.Observable, concurrency int) rx.Operator {
	project := func(interface{}, int) rx.Observable { return inner }

	return MergeSyncMap(project, concurrency)
}

type mergeSyncObservable struct {
	Source      rx.Observable
	Project     func(interface{}, int) rx.Observable
	Concurrency int
}

func (obs mergeSyncObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	sink = sink.WithCancel(cancel).Mutex()

	x := struct {
		sync.Mutex
		Index     int
		Workers   int
		Completed bool
	}{Index: -1}

	var observer rx.Observer

	complete := make(chan struct{}, 1)

	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			x.Lock()

			for x.Workers == obs.Concurrency {
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
