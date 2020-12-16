package operators

import (
	"context"
	"sync"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/queue"
)

// Expand recursively projects each source value to an Observable which is
// merged in the output Observable.
//
// It's similar to MergeMap, but applies the projection function to every
// source value as well as every output value. It's recursive.
//
// For unlimited concurrency, passes -1.
func Expand(
	project func(interface{}) rx.Observable,
	concurrency int,
) rx.Operator {
	if project == nil {
		panic("Expand: project is nil")
	}

	if concurrency == 0 {
		concurrency = -1
	}

	return func(source rx.Observable) rx.Observable {
		return expandObservable{source, project, concurrency}.Subscribe
	}
}

type expandObservable struct {
	Source      rx.Observable
	Project     func(interface{}) rx.Observable
	Concurrency int
}

func (obs expandObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel).Mutex()

	var x struct {
		sync.Mutex
		Queue     queue.Queue
		Workers   int
		Completed bool
	}

	var subscribeLocked func()

	subscribeLocked = func() {
		val := x.Queue.Pop()

		sink.Next(val)

		obs1 := obs.Project(val)

		go obs1.Subscribe(ctx, func(t rx.Notification) {
			switch {
			case t.HasValue:
				x.Lock()
				defer x.Unlock()

				x.Queue.Push(t.Value)

				if x.Workers != obs.Concurrency {
					x.Workers++
					subscribeLocked()
				}

			case t.HasError:
				sink(t)

			default:
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
			}
		})
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			x.Lock()
			defer x.Unlock()

			x.Queue.Push(t.Value)

			if x.Workers != obs.Concurrency {
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
