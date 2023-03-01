package rx

import (
	"context"
)

// SkipWhile skips all values emitted by the source Observable as long as
// a given condition holds true, but emits all further source values as
// soon as the condition becomes false.
func SkipWhile[T any](cond func(v T) bool) Operator[T, T] {
	if cond == nil {
		panic("cond == nil")
	}

	return skipWhile(cond)
}

func skipWhile[T any](cond func(v T) bool) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return skipWhileObservable[T]{source, cond}.Subscribe
		},
	)
}

type skipWhileObservable[T any] struct {
	Source    Observable[T]
	Condition func(T) bool
}

func (obs skipWhileObservable[T]) Subscribe(ctx context.Context, sink Observer[T]) {
	var taking bool

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		switch {
		case taking || !n.HasValue:
			sink(n)
		case !obs.Condition(n.Value):
			taking = true

			sink(n)
		}
	})
}
