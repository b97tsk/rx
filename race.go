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
func RaceWith[T any](some ...Observable[T]) RaceWithOperator[T] {
	return RaceWithOperator[T]{
		opts: raceWithConfig[T]{
			Observables: some,
			PassiveGo:   false,
		},
	}
}

type raceWithConfig[T any] struct {
	Observables []Observable[T]
	PassiveGo   bool
}

// RaceWithOperator is an [Operator] type for [RaceWith].
type RaceWithOperator[T any] struct {
	opts raceWithConfig[T]
}

// WithPassiveGo turns on PassiveGo mode.
// By default, this Operator subscribes to Observables other than the source
// in separate goroutines.
// With PassiveGo mode on, this Operator subscribes to every Observable in
// the same goroutine.
// In PassiveGo mode, goroutines can only be started by Observables themselves.
func (op RaceWithOperator[T]) WithPassiveGo() RaceWithOperator[T] {
	op.opts.PassiveGo = true
	return op
}

// Apply implements the Operator interface.
func (op RaceWithOperator[T]) Apply(source Observable[T]) Observable[T] {
	return raceWithObservable[T]{
		Source:    source,
		Others:    op.opts.Observables,
		PassiveGo: op.opts.PassiveGo,
	}.Subscribe
}

type raceWithObservable[T any] struct {
	Source    Observable[T]
	Others    []Observable[T]
	PassiveGo bool
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
		i += off

		if obs.PassiveGo {
			subscribe(i, obs1)

			if race.Load() != 0 {
				return
			}

			continue
		}

		c.Go(func() { subscribe(i, obs1) })
	}
}

func (obs raceWithObservable[T]) numObservables() int {
	n := len(obs.Others)

	if obs.Source != nil {
		n++
	}

	return n
}
