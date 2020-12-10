package operators

import (
	"context"
	"sync"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/norec"
	"github.com/b97tsk/rx/internal/queue"
)

// ConcatAll creates an Observable that flattens a higher-order Observable into
// a first-order Observable by concatenating the inner Observables in order.
func ConcatAll() rx.Operator {
	return ConcatMap(projectToObservable)
}

// ConcatMap creates an Observable that converts the source Observable into a
// higher-order Observable, by projecting each source value to an Observable,
// and flattens it into a first-order Observable by concatenating the inner
// Observables in order.
func ConcatMap(project func(interface{}, int) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return concatObservable{source, project}.Subscribe
	}
}

// ConcatMapTo creates an Observable that converts the source Observable into
// a higher-order Observable, by projecting each source value to the same
// Observable, and flattens it into a first-order Observable by concatenating
// the inner Observables in order.
//
// It's like ConcatMap, but maps each value always to the same inner Observable.
func ConcatMapTo(inner rx.Observable) rx.Operator {
	return ConcatMap(func(interface{}, int) rx.Observable { return inner })
}

type concatObservable struct {
	Source  rx.Observable
	Project func(interface{}, int) rx.Observable
}

func (obs concatObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel).Mutex()

	x := struct {
		sync.Mutex
		Queue     queue.Queue
		Index     int
		Working   bool
		Completed bool
	}{Index: -1}

	var subscribeToNext func()

	subscribeToNext = norec.Wrap(func() {
		x.Lock()

		if x.Queue.Len() == 0 {
			defer x.Unlock()

			x.Working = false

			if x.Completed {
				sink.Complete()
			}

			return
		}

		x.Index++

		sourceIndex := x.Index
		sourceValue := x.Queue.Pop()

		x.Unlock()

		obs1 := obs.Project(sourceValue, sourceIndex)

		obs1.Subscribe(ctx, func(t rx.Notification) {
			if t.HasValue || t.HasError {
				sink(t)
				return
			}

			if ctx.Err() != nil {
				return
			}

			subscribeToNext()
		})
	})

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			x.Lock()

			x.Queue.Push(t.Value)

			var subscribe bool

			if !x.Working {
				x.Working = true
				subscribe = true
			}

			x.Unlock()

			if subscribe {
				subscribeToNext()
			}

		case t.HasError:
			sink(t)

		default:
			x.Lock()
			defer x.Unlock()

			x.Completed = true

			if !x.Working {
				sink(t)
			}
		}
	})
}
