package rx

import (
	"context"
)

// Contains emits a boolean to indicate whether any value of the source
// Observable satisfies a given condition.
func Contains[T any](cond func(v T) bool) Operator[T, bool] {
	if cond == nil {
		panic("cond == nil")
	}

	return contains(cond)
}

func contains[T any](cond func(v T) bool) Operator[T, bool] {
	return NewOperator(
		func(source Observable[T]) Observable[bool] {
			return containsObservable[T]{source, cond}.Subscribe
		},
	)
}

// ContainsElement emits a boolean to indicate whether the source Observable
// emits a given value.
func ContainsElement[T comparable](v T) Operator[T, bool] {
	return NewOperator(
		func(source Observable[T]) Observable[bool] {
			return containsElementObservable[T]{source, v}.Subscribe
		},
	)
}

type containsObservable[T any] struct {
	Source    Observable[T]
	Condition func(T) bool
}

func (obs containsObservable[T]) Subscribe(ctx context.Context, sink Observer[bool]) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancel)

	var noop bool

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		switch {
		case noop:
		case n.HasValue:
			if obs.Condition(n.Value) {
				noop = true

				sink.Next(true)
				sink.Complete()
			}
		case n.HasError:
			sink.Error(n.Error)
		default:
			sink.Next(false)
			sink.Complete()
		}
	})
}

type containsElementObservable[T comparable] struct {
	Source  Observable[T]
	Element T
}

func (obs containsElementObservable[T]) Subscribe(ctx context.Context, sink Observer[bool]) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancel)

	var noop bool

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		switch {
		case noop:
		case n.HasValue:
			if n.Value == obs.Element {
				noop = true

				sink.Next(true)
				sink.Complete()
			}
		case n.HasError:
			sink.Error(n.Error)
		default:
			sink.Next(false)
			sink.Complete()
		}
	})
}
