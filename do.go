package rx

import (
	"context"
)

// Do creates an Observable that mirrors the source Observable, but performs
// a side effect before each emission.
func (Operators) Do(sink Observer) Operator {
	return func(source Observable) Observable {
		return func(ctx context.Context, notify Observer) (context.Context, context.CancelFunc) {
			return source.Subscribe(ctx, func(t Notification) {
				sink(t)
				notify(t)
			})
		}
	}
}

// DoOnNext creates an Observable that mirrors the source Observable, but
// performs a side effect before each value.
func (Operators) DoOnNext(onNext func(interface{})) Operator {
	return func(source Observable) Observable {
		return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
			return source.Subscribe(ctx, func(t Notification) {
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
func (Operators) DoOnError(onError func(error)) Operator {
	return func(source Observable) Observable {
		return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
			return source.Subscribe(ctx, func(t Notification) {
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
func (Operators) DoOnComplete(onComplete func()) Operator {
	return func(source Observable) Observable {
		return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
			return source.Subscribe(ctx, func(t Notification) {
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

// DoAtLast creates an Observer that passes all emissions to the specified
// Observer, in the case that an ERROR or COMPLETE emission is passed, makes
// a call to the specified function.
func DoAtLast(sink Observer, atLast func(Notification)) Observer {
	return func(t Notification) {
		if t.HasValue {
			sink(t)
			return
		}
		sink(t)
		atLast(t)
	}
}

// DoAtLast creates an Observable that mirrors the source Observable, in the
// case that an ERROR or COMPLETE emission is mirrored, makes a call to the
// specified function.
func (Operators) DoAtLast(atLast func(Notification)) Operator {
	return func(source Observable) Observable {
		return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
			return source.Subscribe(ctx, DoAtLast(sink, atLast))
		}
	}
}
