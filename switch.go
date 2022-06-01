package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/atomic"
	"github.com/b97tsk/rx/internal/critical"
)

// SwitchAll flattens a higher-order Observable into a first-order Observable
// by subscribing to only the most recently emitted of those inner Observables.
func SwitchAll[_ Observable[T], T any]() Operator[Observable[T], T] {
	return switchMap(identity[Observable[T]])
}

// SwitchMap converts the source Observable into a higher-order Observable,
// by projecting each source value to an Observable, then flattens it into
// a first-order Observable using SwitchAll.
func SwitchMap[T, R any](proj func(v T) Observable[R]) Operator[T, R] {
	if proj == nil {
		panic("proj == nil")
	}

	return switchMap(proj)
}

// SwitchMapTo converts the source Observable into a higher-order Observable,
// by projecting each source value to the same Observable, then flattens it
// into a first-order Observable using SwitchAll.
func SwitchMapTo[T, R any](inner Observable[R]) Operator[T, R] {
	if inner == nil {
		panic("inner == nil")
	}

	return switchMap(func(T) Observable[R] { return inner })
}

func switchMap[T, R any](proj func(v T) Observable[R]) Operator[T, R] {
	return AsOperator(
		func(source Observable[T]) Observable[R] {
			return switchMapObservable[T, R]{source, proj}.Subscribe
		},
	)
}

type switchMapObservable[T, R any] struct {
	Source  Observable[T]
	Project func(T) Observable[R]
}

func (obs switchMapObservable[T, R]) Subscribe(ctx context.Context, sink Observer[R]) {
	ctx, cancel := context.WithCancel(ctx)

	var x struct {
		critical.Section
		Ctx       atomic.Value
		Cancel    context.CancelFunc
		Completed atomic.Bool
	}

	x.Ctx.Store(ctx)

	sinkAndDone := func(n Notification[R]) {
		if critical.Enter(&x.Section) {
			critical.Close(&x.Section)
			cancel()
			sink(n)
		}
	}

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		switch {
		case n.HasValue:
			childCtx, cancel := context.WithCancel(ctx)
			x.Ctx.Store(childCtx)

			if x.Cancel != nil {
				x.Cancel()
			}

			x.Cancel = cancel

			obs.Project(n.Value).Subscribe(childCtx, func(n Notification[R]) {
				if n.HasValue {
					if critical.Enter(&x.Section) {
						if x.Ctx.Load() == childCtx {
							sink(n)
						}

						critical.Leave(&x.Section)
					}

					return
				}

				cancel()

				if x.Ctx.CompareAndSwap(childCtx, ctx) {
					if n.HasError || x.Completed.True() && x.Ctx.Load() == ctx {
						sinkAndDone(n)
					}
				}
			})

		case n.HasError:
			sinkAndDone(Error[R](n.Error))

		default:
			x.Completed.Store(true)

			if x.Ctx.Load() == ctx {
				sinkAndDone(Complete[R]())
			}
		}
	})
}
