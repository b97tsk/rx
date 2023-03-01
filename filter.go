package rx

import (
	"context"
)

// Filter filters values emitted by the source Observable by only emitting
// those that satisfy a given condition.
func Filter[T any](cond func(v T) bool) Operator[T, T] {
	if cond == nil {
		panic("cond == nil")
	}

	return filter(cond)
}

func filter[T any](cond func(v T) bool) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				source.Subscribe(ctx, func(n Notification[T]) {
					switch {
					case n.HasValue:
						if cond(n.Value) {
							sink(n)
						}
					default:
						sink(n)
					}
				})
			}
		},
	)
}

// FilterOut filters out values emitted by the source Observable by only
// emitting those that do not satisfy a given condition.
func FilterOut[T any](cond func(v T) bool) Operator[T, T] {
	if cond == nil {
		panic("cond == nil")
	}

	return filterOut(cond)
}

func filterOut[T any](cond func(v T) bool) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				source.Subscribe(ctx, func(n Notification[T]) {
					switch {
					case n.HasValue:
						if !cond(n.Value) {
							sink(n)
						}
					default:
						sink(n)
					}
				})
			}
		},
	)
}

// FilterMap passes each value emitted by the source Observable to a given
// condition function and emits their mapping, the first return value of
// the condition function, only if the second is true.
func FilterMap[T, R any](cond func(v T) (R, bool)) Operator[T, R] {
	if cond == nil {
		panic("cond == nil")
	}

	return filterMap(cond)
}

func filterMap[T, R any](cond func(v T) (R, bool)) Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return func(ctx context.Context, sink Observer[R]) {
				source.Subscribe(ctx, func(n Notification[T]) {
					switch {
					case n.HasValue:
						if v, ok := cond(n.Value); ok {
							sink.Next(v)
						}
					case n.HasError:
						sink.Error(n.Error)
					default:
						sink.Complete()
					}
				})
			}
		},
	)
}
