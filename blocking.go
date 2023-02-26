package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/waitgroup"
)

// BlockingFirst subscribes to the source Observable, returns the first item
// emitted by the source; if the source emits no items, it returns zero value
// of T and ErrEmpty; if the source throws an error, it returns zero value of
// T and this error.
//
// A cancellation of ctx will cause BlockingFirst to immediately return zero
// value of T and ctx.Err().
func (obs Observable[T]) BlockingFirst(ctx context.Context) (v T, err error) {
	childCtx, cancel := context.WithCancel(ctx)

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
	})

	<-childCtx.Done()

	if err := getErr(ctx); err != nil {
		return v, err
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

// BlockingFirstOrDefault subscribes to the source Observable, returns
// the first item emitted by the source, or returns def if the source emits
// no items or throws an error.
//
// A cancellation of ctx will cause BlockingFirstOrDefault to immediately
// return def.
func (obs Observable[T]) BlockingFirstOrDefault(ctx context.Context, def T) T {
	v, err := obs.BlockingFirst(ctx)
	if err != nil {
		return def
	}

	return v
}

// BlockingLast subscribes to the source Observable, returns the last item
// emitted by the source; if the source emits no items, it returns zero value
// of T and ErrEmpty; if the source throws an error, it returns zero value of
// T and this error.
//
// A cancellation of ctx will cause BlockingLast to immediately return zero
// value of T and ctx.Err().
func (obs Observable[T]) BlockingLast(ctx context.Context) (v T, err error) {
	childCtx, cancel := context.WithCancel(ctx)

	res := Error[T](ErrEmpty)

	obs.Subscribe(childCtx, func(n Notification[T]) {
		if n.HasValue || n.HasError {
			res = n
		}

		if !n.HasValue {
			cancel()
		}
	})

	<-childCtx.Done()

	if err := getErr(ctx); err != nil {
		return v, err
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

// BlockingLastOrDefault subscribes to the source Observable, returns the last
// item emitted by the source, or returns def if the source emits no items or
// throws an error.
//
// A cancellation of ctx will cause BlockingLastOrDefault to immediately return
// def.
func (obs Observable[T]) BlockingLastOrDefault(ctx context.Context, def T) T {
	v, err := obs.BlockingLast(ctx)
	if err != nil {
		return def
	}

	return v
}

// BlockingSingle subscribes to the source Observable, returns the single item
// emitted by the source; if the source emits more than one item or no items,
// it returns zero value of T and ErrNotSingle or ErrEmpty respectively;
// if the source throws an error, it returns zero value of T and this error.
//
// A cancellation of ctx will cause BlockingSingle to immediately return zero
// value of T and ctx.Err().
func (obs Observable[T]) BlockingSingle(ctx context.Context) (v T, err error) {
	childCtx, cancel := context.WithCancel(ctx)

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

			return
		}

		if n.HasValue || n.HasError {
			res = n
		}

		if !n.HasValue {
			cancel()
		}
	})

	<-childCtx.Done()

	if err := getErr(ctx); err != nil {
		return v, err
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

// BlockingSubscribe subscribes to the source Observable, returns only when
// the source completes or throws an error; if the source completes, it returns
// nil; if the source throws an error, it returns this error.
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
