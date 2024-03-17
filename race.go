package rx

import "sync/atomic"

// Race creates an Observable that mirrors the first Observable to emit
// a value, from given input Observables.
func Race[T any](some ...Observable[T]) Observable[T] {
	if len(some) == 0 {
		return Empty[T]()
	}

	return raceWithObservable[T]{Others: some}.Subscribe
}

// RaceWith mirrors the first Observable to emit a value, from the source
// and given input Observables.
func RaceWith[T any](some ...Observable[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return raceWithObservable[T]{source, some}.Subscribe
		},
	)
}

type raceWithObservable[T any] struct {
	Source Observable[T]
	Others []Observable[T]
}

func (obs raceWithObservable[T]) Subscribe(c Context, sink Observer[T]) {
	subs := make([]Pair[Context, CancelFunc], obs.numObservables())

	for i := range subs {
		subs[i] = NewPair(c.WithCancel())
	}

	var race atomic.Uint32

	subscribe := func(i int, obs Observable[T]) {
		var won, lost bool

		obs.Subscribe(subs[i].Left(), func(n Notification[T]) {
			switch {
			case won:
				sink(n)
				return
			case lost:
				return
			}

			if race.CompareAndSwap(0, 1) {
				for j := range subs {
					if j != i {
						subs[j].Right()()
					}
				}

				won = true

				sink(n)

				return
			}

			lost = true

			subs[i].Right()()
		})
	}

	var off int

	if obs.Source != nil {
		subscribe(0, obs.Source)

		if race.Load() != 0 {
			return
		}

		off = 1
	}

	for i, obs1 := range obs.Others {
		subscribe(i+off, obs1)

		if race.Load() != 0 {
			return
		}
	}
}

func (obs raceWithObservable[T]) numObservables() int {
	n := len(obs.Others)

	if obs.Source != nil {
		n++
	}

	return n
}
