package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Do mirrors the source, but performs a side effect before each emission.
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

// DoOnNext mirrors the source, but performs a side effect before each value.
func DoOnNext(f func(interface{})) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			source.Subscribe(ctx, func(t rx.Notification) {
				if t.HasValue {
					f(t.Value)
				}

				sink(t)
			})
		}
	}
}

// DoOnError mirrors the source and, when the source throws an error, performs
// a side effect before mirroring this error.
func DoOnError(f func(error)) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			source.Subscribe(ctx, func(t rx.Notification) {
				if t.HasError {
					f(t.Error)
				}

				sink(t)
			})
		}
	}
}

// DoOnComplete mirrors the source and, when the source completes, performs
// a side effect before mirroring this completion.
func DoOnComplete(f func()) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			source.Subscribe(ctx, func(t rx.Notification) {
				if !t.HasValue && !t.HasError {
					f()
				}

				sink(t)
			})
		}
	}
}

// DoOnErrorOrComplete mirrors the source and, when the source throws an
// error or completes, performs a side effect before mirroring this error
// or completion.
func DoOnErrorOrComplete(f func()) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			source.Subscribe(ctx, func(t rx.Notification) {
				if !t.HasValue {
					f()
				}

				sink(t)
			})
		}
	}
}
