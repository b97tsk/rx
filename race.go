package rx

import (
	"context"
	"sync/atomic"
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

	var race atomic.Uint32

	wg := WaitGroupFromContext(ctx)

	for index, obs := range some {
		index, obs := index, obs

		wg.Go(func() {
			var won, lost bool

			obs.Subscribe(subs[index].Left(), func(n Notification[T]) {
				switch {
				case won:
					sink(n)
					return
				case lost:
					return
				}

				if race.CompareAndSwap(0, 1) {
					for i := range subs {
						if i != index {
							subs[i].Right()()
						}
					}

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
