package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/atomic"
)

// SkipUntil skips items emitted by the source Observable until a second
// Observable emits an item.
func SkipUntil[T, U any](notifier Observable[U]) Operator[T, T] {
	if notifier == nil {
		panic("notifier == nil")
	}

	return skipUntil[T](notifier)
}

func skipUntil[T, U any](notifier Observable[U]) Operator[T, T] {
	return AsOperator(
		func(source Observable[T]) Observable[T] {
			return skipUntilObservable[T, U]{source, notifier}.Subscribe
		},
	)
}

type skipUntilObservable[T, U any] struct {
	Source   Observable[T]
	Notifier Observable[U]
}

func (obs skipUntilObservable[T, U]) Subscribe(ctx context.Context, sink Observer[T]) {
	ctx, cancel := context.WithCancel(ctx)

	originalSink := sink

	sink = sink.WithCancel(cancel).Mutex()

	var noSkipping atomic.Bool

	{
		ctx, cancel := context.WithCancel(ctx)

		var noop bool

		obs.Notifier.Subscribe(ctx, func(n Notification[U]) {
			if noop {
				return
			}

			noop = true

			cancel()

			switch {
			case n.HasValue:
				noSkipping.Store(true)
			case n.HasError:
				sink.Error(n.Error)
			}
		})
	}

	if err := ctx.Err(); err != nil {
		sink.Error(err)
		return
	}

	if noSkipping.True() {
		obs.Source.Subscribe(ctx, originalSink)
		return
	}

	var taking bool

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		switch {
		case taking:
			originalSink(n)
		case n.HasValue:
			if noSkipping.True() {
				taking = true

				originalSink(n)
			}
		default:
			sink(n)
		}
	})
}