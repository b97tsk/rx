package rx

import (
	"context"
)

// Every emits a boolean to indicate whether every value of the source
// Observable satisfies a given condition.
func Every[T any](cond func(v T) bool) Operator[T, bool] {
	if cond == nil {
		panic("cond == nil")
	}

	return every(cond)
}

func every[T any](cond func(v T) bool) Operator[T, bool] {
	return NewOperator(
		func(source Observable[T]) Observable[bool] {
			return everyObservable[T]{source, cond}.Subscribe
		},
	)
}

type everyObservable[T any] struct {
	Source    Observable[T]
	Condition func(T) bool
}

func (obs everyObservable[T]) Subscribe(ctx context.Context, sink Observer[bool]) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var noop bool

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		switch {
		case noop:
		case n.HasValue:
			if !obs.Condition(n.Value) {
				noop = true

				sink.Next(false)
				sink.Complete()
			}
		case n.HasError:
			sink.Error(n.Error)
		default:
			sink.Next(true)
			sink.Complete()
		}
	})
}
