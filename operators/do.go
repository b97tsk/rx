package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Do creates an Observable that mirrors the source Observable, but performs
// a side effect before each emission.
func Do(tap rx.Observer) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			source.Subscribe(ctx, func(t rx.Notification) {
				tap(t)
				sink(t)
			})
		}
	}
}

// DoOnNext creates an Observable that mirrors the source Observable, but
// performs a side effect before each value.
func DoOnNext(onNext func(interface{})) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			source.Subscribe(ctx, func(t rx.Notification) {
				if t.HasValue {
					onNext(t.Value)
				}

				sink(t)
			})
		}
	}
}

// DoOnError creates an Observable that mirrors the source Observable and,
// when the source emits an error, performs a side effect before mirroring
// this error.
func DoOnError(onError func(error)) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			source.Subscribe(ctx, func(t rx.Notification) {
				if t.HasError {
					onError(t.Error)
				}

				sink(t)
			})
		}
	}
}

// DoOnComplete creates an Observable that mirrors the source Observable
// and, when the source completes, performs a side effect before mirroring
// this completion.
func DoOnComplete(onComplete func()) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			source.Subscribe(ctx, func(t rx.Notification) {
				if !t.HasValue && !t.HasError {
					onComplete()
				}

				sink(t)
			})
		}
	}
}

// DoAtLast creates an Observable that mirrors the source Observable and,
// when the source completes or emits an error, performs a side effect after.
//
// If you also want to know whether there's an error, use DoAtLastError
// instead.
func DoAtLast(atLast func()) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			source.Subscribe(ctx, func(t rx.Notification) {
				if t.HasValue {
					sink(t)

					return
				}

				sink(t)
				atLast()
			})
		}
	}
}

// DoAtLastError creates an Observable that mirrors the source Observable and,
// when the source completes or emits an error, performs a side effect after.
//
// If you need not to know whether there's an error, use DoAtLast instead.
func DoAtLastError(atLast func(error)) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			source.Subscribe(ctx, func(t rx.Notification) {
				if t.HasValue {
					sink(t)

					return
				}

				var e error

				if t.HasError {
					e = t.Error
				}

				sink(t)
				atLast(e)
			})
		}
	}
}
