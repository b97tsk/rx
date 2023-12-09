package rx

import (
	"context"
)

// TakeWhile emits values emitted by the source Observable so long as
// each value satisfies a given condition, and then completes as soon as
// the condition is not satisfied.
func TakeWhile[T any](cond func(v T) bool) Operator[T, T] {
	if cond == nil {
		panic("cond == nil")
	}

	return takeWhile(cond)
}

func takeWhile[T any](cond func(v T) bool) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return takeWhileObservable[T]{source, cond}.Subscribe
		},
	)
}

type takeWhileObservable[T any] struct {
	Source    Observable[T]
	Condition func(T) bool
}

func (obs takeWhileObservable[T]) Subscribe(ctx context.Context, sink Observer[T]) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancel)

	var noop bool

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		if noop {
			return
		}

		switch n.Kind {
		case KindNext:
			if obs.Condition(n.Value) {
				sink(n)
				return
			}

			noop = true

			sink.Complete()

		case KindError, KindComplete:
			sink(n)
		}
	})
}
