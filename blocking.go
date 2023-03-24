package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/waitgroup"
)

// BlockingFirst subscribes to the source Observable, and returns
// the first value emitted by the source.
// If the source emits no values, it returns zero value of T and ErrEmpty;
// if the source emits a notification of error, it returns zero value of T
// and the error.
//
// The source Observable must honor the cancellation of ctx; otherwise,
// BlockingFirst might still block even after ctx has been cancelled.
func (obs Observable[T]) BlockingFirst(ctx context.Context) (v T, err error) {
	childCtx, cancel := context.WithCancel(ctx)
	childCtx, wg := waitgroup.Install(childCtx)

	res := Error[T](ErrEmpty)

	var noop bool

	obs.Subscribe(childCtx, func(n Notification[T]) {
		if noop {
			return
		}

		noop = true

		if n.HasValue || n.HasError {
			res = n
		}

		cancel()
		wg.Done()
	})

	wg.Wait()

	select {
	default:
	case <-ctx.Done():
		return v, ctx.Err()
	}

	switch {
	case res.HasValue:
		return res.Value, nil
	case res.HasError:
		return v, res.Error
	default:
		panic("unreachable")
	}
}

// BlockingFirstOrElse subscribes to the source Observable, and returns
// the first value emitted by the source, or returns def if the source emits
// no values or a notification of error.
//
// The source Observable must honor the cancellation of ctx; otherwise,
// BlockingFirstOrElse might still block even after ctx has been cancelled.
func (obs Observable[T]) BlockingFirstOrElse(ctx context.Context, def T) T {
	v, err := obs.BlockingFirst(ctx)
	if err != nil {
		return def
	}

	return v
}

// BlockingLast subscribes to the source Observable, and returns
// the last value emitted by the source.
// If the source emits no values, it returns zero value of T and ErrEmpty;
// if the source emits a notification of error, it returns zero value of T
// and the error.
//
// The source Observable must honor the cancellation of ctx; otherwise,
// BlockingLast might still block even after ctx has been cancelled.
func (obs Observable[T]) BlockingLast(ctx context.Context) (v T, err error) {
	childCtx, cancel := context.WithCancel(ctx)
	childCtx, wg := waitgroup.Install(childCtx)

	res := Error[T](ErrEmpty)

	obs.Subscribe(childCtx, func(n Notification[T]) {
		if n.HasValue || n.HasError {
			res = n
		}

		if !n.HasValue {
			cancel()
			wg.Done()
		}
	})

	wg.Wait()

	select {
	default:
	case <-ctx.Done():
		return v, ctx.Err()
	}

	switch {
	case res.HasValue:
		return res.Value, nil
	case res.HasError:
		return v, res.Error
	default:
		panic("unreachable")
	}
}

// BlockingLastOrElse subscribes to the source Observable, and returns
// the last value emitted by the source, or returns def if the source emits
// no values or a notification of error.
//
// The source Observable must honor the cancellation of ctx; otherwise,
// BlockingLastOrElse might still block even after ctx has been cancelled.
func (obs Observable[T]) BlockingLastOrElse(ctx context.Context, def T) T {
	v, err := obs.BlockingLast(ctx)
	if err != nil {
		return def
	}

	return v
}

// BlockingSingle subscribes to the source Observable, and returns
// the single value emitted by the source.
// If the source emits more than one value or no values, it returns
// zero value of T and ErrNotSingle or ErrEmpty respectively;
// if the source emits a notification of error, it returns
// zero value of T and the error.
//
// The source Observable must honor the cancellation of ctx; otherwise,
// BlockingSingle might still block even after ctx has been cancelled.
func (obs Observable[T]) BlockingSingle(ctx context.Context) (v T, err error) {
	childCtx, cancel := context.WithCancel(ctx)
	childCtx, wg := waitgroup.Install(childCtx)

	res := Error[T](ErrEmpty)

	var noop bool

	obs.Subscribe(childCtx, func(n Notification[T]) {
		if noop {
			return
		}

		if n.HasValue && res.HasValue {
			res = Error[T](ErrNotSingle)
			noop = true

			cancel()
			wg.Done()

			return
		}

		if n.HasValue || n.HasError {
			res = n
		}

		if !n.HasValue {
			cancel()
			wg.Done()
		}
	})

	wg.Wait()

	select {
	default:
	case <-ctx.Done():
		return v, ctx.Err()
	}

	switch {
	case res.HasValue:
		return res.Value, nil
	case res.HasError:
		return v, res.Error
	default:
		panic("unreachable")
	}
}

// BlockingSubscribe subscribes to the source Observable, and returns nil
// when the source completes, or returns an error when the source emits
// a notification of error.
//
// The source Observable must honor the cancellation of ctx; otherwise,
// BlockingSubscribe might still block even after ctx has been cancelled.
func (obs Observable[T]) BlockingSubscribe(ctx context.Context, sink Observer[T]) error {
	var res Notification[T]

	ctx, wg := waitgroup.Install(ctx)

	obs.Subscribe(ctx, func(n Notification[T]) {
		res = n

		sink(n)

		if !n.HasValue {
			wg.Done()
		}
	})

	wg.Wait()

	switch {
	case res.HasValue:
		panic("unreachable")
	case res.HasError:
		return res.Error
	default:
		return nil
	}
}
