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
func Expand(project func(interface{}) rx.Observable) rx.Operator {
	return ExpandConfigure{project, -1}.Make()
}

// An ExpandConfigure is a configure for Expand.
type ExpandConfigure struct {
	Project     func(interface{}) rx.Observable
	Concurrency int
}

// Make creates an Operator from this configure.
func (configure ExpandConfigure) Make() rx.Operator {
	if configure.Project == nil {
		panic("Expand: Project is nil")
	}
	if configure.Concurrency == 0 {
		configure.Concurrency = -1
	}
	return func(source rx.Observable) rx.Observable {
		return expandObservable{source, configure}.Subscribe
	}
}

type expandObservable struct {
	Source rx.Observable
	ExpandConfigure
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
