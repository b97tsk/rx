package rx

import (
	"context"
)

// Find emits only the first value emitted by the source Observable that
// meets some condition.
func Find[T any](cond func(v T) bool) Operator[T, T] {
	if cond == nil {
		panic("cond == nil")
	}

	return find(cond)
}

func find[T any](cond func(v T) bool) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return findObservable[T]{source, cond}.Subscribe
		},
	)
}

type findObservable[T any] struct {
	Source    Observable[T]
	Condition func(T) bool
}

func (obs findObservable[T]) Subscribe(ctx context.Context, sink Observer[T]) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancel)

	var noop bool

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		switch {
		case noop:
		case n.HasValue:
			if obs.Condition(n.Value) {
				noop = true

				sink(n)
				sink.Complete()
			}
		default:
			sink(n)
		}
	})
}
