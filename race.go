package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/critical"
	"github.com/b97tsk/rx/internal/waitgroup"
)

// Race creates an Observable that mirrors the first Observable to emit
// a value, from given input Observables.
func Race[T any](some ...Observable[T]) Observable[T] {
	if len(some) == 0 {
		return Empty[T]()
	}

	return observables[T](some).Race
}

// RaceWith applies [Race] to the source Observable along with some other
// Observables to create a first-order Observable.
func RaceWith[T any](some ...Observable[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return observables[T](append([]Observable[T]{source}, some...)).Race
		},
	)
}

func (some observables[T]) Race(ctx context.Context, sink Observer[T]) {
	subs := make([]Pair[context.Context, context.CancelFunc], len(some))

	for i := range subs {
		subs[i] = NewPair(context.WithCancel(ctx))
	}

	var race critical.Section

	ctxHoisted := waitgroup.Hoist(ctx)

	for index, obs := range some {
		index, obs := index, obs

		var won, lost bool

		Go(ctxHoisted, func() {
			obs.Subscribe(subs[index].Left(), func(n Notification[T]) {
				switch {
				case won:
					sink(n)
					return
				case lost:
					return
				}

				if critical.Enter(&race) {
					for i := range subs {
						if i != index {
							subs[i].Right()()
						}
					}

					critical.Close(&race)

					won = true

					sink(n)

					return
				}

				lost = true

				subs[index].Right()()
			})
		})
	}
}
