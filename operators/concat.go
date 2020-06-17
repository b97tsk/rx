package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/misc"
	"github.com/b97tsk/rx/internal/queue"
)

type concatObservable struct {
	Source  rx.Observable
	Project func(interface{}, int) rx.Observable
}

func (obs concatObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	sink = rx.Mutex(sink)

	type X struct {
		Index  int
		Active int
		Buffer queue.Queue
	}
	cx := make(chan *X, 1)
	cx <- &X{Active: 1}

	var doNextLocked func(*X)

	doNextLocked = func(x *X) {
		var avoidRecursion misc.AvoidRecursion
		avoidRecursion.Do(func() {
			if x.Buffer.Len() == 0 {
				x.Active--
				if x.Active == 0 {
					sink.Complete()
				}
				return
			}

			sourceIndex := x.Index
			sourceValue := x.Buffer.Pop()
			x.Index++

			obs1 := obs.Project(sourceValue, sourceIndex)

			obs1.Subscribe(ctx, func(t rx.Notification) {
				if t.HasValue || t.HasError {
					sink(t)
					return
				}
				if ctx.Err() != nil {
					return
				}
				avoidRecursion.Do(func() {
					x := <-cx
					doNextLocked(x)
					cx <- x
				})
			})
		})
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			x := <-cx
			x.Buffer.Push(t.Value)
			if x.Active == 1 {
				x.Active++
				doNextLocked(x)
			}
			cx <- x

		case t.HasError:
			sink(t)

		default:
			x := <-cx
			x.Active--
			if x.Active == 0 {
				sink(t)
			}
			cx <- x
		}
	})
}

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
		obs := concatObservable{source, project}
		return rx.Create(obs.Subscribe)
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
