package rx

import (
	"context"
)

// Do creates an Observable that mirrors the source Observable, but performs
// a side effect before each emission.
func (Operators) Do(sink Observer) OperatorFunc {
	return func(source Observable) Observable {
		return source.Lift(
			func(ctx context.Context, notify Observer, source Observable) (context.Context, context.CancelFunc) {
				return source.Subscribe(ctx, func(t Notification) {
					sink(t)
					notify(t)
				})
			},
		)
	}
}

// DoOnNext creates an Observable that mirrors the source Observable, but
// performs a side effect before each value.
func (Operators) DoOnNext(onNext func(interface{})) OperatorFunc {
	return func(source Observable) Observable {
		return source.Lift(
			func(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
				return source.Subscribe(ctx, func(t Notification) {
					if t.HasValue {
						onNext(t.Value)
					}
					sink(t)
				})
			},
		)
	}
}

// DoOnError creates an Observable that mirrors the source Observable, in
// the case that the source emits an Error emission, performs a side effect
// before mirroring that emission.
func (Operators) DoOnError(onError func(error)) OperatorFunc {
	return func(source Observable) Observable {
		return source.Lift(
			func(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
				return source.Subscribe(ctx, func(t Notification) {
					if t.HasError {
						onError(t.Value.(error))
					}
					sink(t)
				})
			},
		)
	}
}

// DoOnComplete creates an Observable that mirrors the source Observable, in
// the case that the source emits an Error emission, performs a side effect
// before mirroring that emission.
func (Operators) DoOnComplete(onComplete func()) OperatorFunc {
	return func(source Observable) Observable {
		return source.Lift(
			func(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
				return source.Subscribe(ctx, func(t Notification) {
					if t.HasValue || t.HasError {
						sink(t)
						return
					}
					onComplete()
					sink(t)
				})
			},
		)
	}
}
