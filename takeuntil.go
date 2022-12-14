package rx

import (
	"context"
)

// TakeUntil mirrors the source Observable until a second Observable emits
// a value.
func TakeUntil[T, U any](notifier Observable[U]) Operator[T, T] {
	if notifier == nil {
		panic("notifier == nil")
	}

	return takeUntil[T](notifier)
}

func takeUntil[T, U any](notifier Observable[U]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return takeUntilObservable[T, U]{source, notifier}.Subscribe
		},
	)
}

type takeUntilObservable[T, U any] struct {
	Source   Observable[T]
	Notifier Observable[U]
}

func (obs takeUntilObservable[T, U]) Subscribe(ctx context.Context, sink Observer[T]) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel).Mutex()

	obs.Notifier.Subscribe(ctx, func(n Notification[U]) {
		switch {
		case n.HasValue:
			sink.Complete()
		case n.HasError:
			sink.Error(n.Error)
		}
	})

	if err := getErr(ctx); err != nil {
		sink.Error(err)
		return
	}

	obs.Source.Subscribe(ctx, sink)
}
