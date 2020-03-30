package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Do creates an Observable that mirrors the source Observable, but performs
// a side effect before each emission.
func Do(sink rx.Observer) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, notify rx.Observer) (context.Context, context.CancelFunc) {
			return source.Subscribe(ctx, func(t rx.Notification) {
				sink(t)
				notify(t)
			})
		}
	}
}

// DoOnNext creates an Observable that mirrors the source Observable, but
// performs a side effect before each value.
func DoOnNext(onNext func(interface{})) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) (context.Context, context.CancelFunc) {
			return source.Subscribe(ctx, func(t rx.Notification) {
				if t.HasValue {
					onNext(t.Value)
				}
				sink(t)
			})
		}
	}
}

// DoOnError creates an Observable that mirrors the source Observable, in
// the case that the source errors, performs a side effect before mirroring
// the ERROR emission.
func DoOnError(onError func(error)) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) (context.Context, context.CancelFunc) {
			return source.Subscribe(ctx, func(t rx.Notification) {
				if t.HasError {
					onError(t.Error)
				}
				sink(t)
			})
		}
	}
}

// DoOnComplete creates an Observable that mirrors the source Observable, in
// the case that the source completes, performs a side effect before mirroring
// the COMPLETE emission.
func DoOnComplete(onComplete func()) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) (context.Context, context.CancelFunc) {
			return source.Subscribe(ctx, func(t rx.Notification) {
				if t.HasValue || t.HasError {
					sink(t)
					return
				}
				onComplete()
				sink(t)
			})
		}
	}
}

// DoAtLast creates an Observable that mirrors the source Observable, in the
// case that an ERROR or COMPLETE emission is mirrored, makes a call to the
// specified function.
func DoAtLast(atLast func(rx.Notification)) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) (context.Context, context.CancelFunc) {
			return source.Subscribe(ctx, func(t rx.Notification) {
				if t.HasValue {
					sink(t)
					return
				}
				sink(t)
				atLast(t)
			})
		}
	}
}
