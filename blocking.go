package rx

import "context"

// BlockingFirst subscribes to the source Observable, and returns
// the first value emitted by the source.
// If the source emits no values, it returns zero value of T and ErrEmpty;
// if the source emits a notification of error, it returns zero value of T
// and the error.
//
// The source Observable must honor the cancellation of ctx; otherwise,
// BlockingFirst might still block even after ctx has been cancelled.
//
// Like any other Blocking methods, this method waits for every goroutine
// started during subscription to complete before returning.
// To have this work properly, Observables must use [WaitGroupFromContext]
// to obtain a WaitGroup and use [WaitGroup.Go] rather than built-in go
// statements to start new goroutines during subscription, especially when
// they need to subscribe to other Observables in a goroutine; otherwise,
// runtime panicking might happen randomly (WaitGroup misuse).
func (obs Observable[T]) BlockingFirst(ctx context.Context) (v T, err error) {
	var wg WaitGroup

	child, cancel := context.WithCancel(WithWaitGroup(ctx, &wg))

	res := Error[T](ErrEmpty)

	var noop bool

	wg.Add(1)

	obs.Subscribe(child, func(n Notification[T]) {
		if noop {
			return
		}

		switch n.Kind {
		case KindNext, KindError, KindComplete:
			noop = true

			switch n.Kind {
			case KindNext, KindError:
				res = n
			}

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

	switch res.Kind {
	case KindNext:
		return res.Value, nil
	case KindError:
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
//
// Like any other Blocking methods, this method waits for every goroutine
// started during subscription to complete before returning.
// To have this work properly, Observables must use [WaitGroupFromContext]
// to obtain a WaitGroup and use [WaitGroup.Go] rather than built-in go
// statements to start new goroutines during subscription, especially when
// they need to subscribe to other Observables in a goroutine; otherwise,
// runtime panicking might happen randomly (WaitGroup misuse).
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
//
// Like any other Blocking methods, this method waits for every goroutine
// started during subscription to complete before returning.
// To have this work properly, Observables must use [WaitGroupFromContext]
// to obtain a WaitGroup and use [WaitGroup.Go] rather than built-in go
// statements to start new goroutines during subscription, especially when
// they need to subscribe to other Observables in a goroutine; otherwise,
// runtime panicking might happen randomly (WaitGroup misuse).
func (obs Observable[T]) BlockingLast(ctx context.Context) (v T, err error) {
	var wg WaitGroup

	child, cancel := context.WithCancel(WithWaitGroup(ctx, &wg))

	res := Error[T](ErrEmpty)

	wg.Add(1)

	obs.Subscribe(child, func(n Notification[T]) {
		switch n.Kind {
		case KindNext, KindError:
			res = n
		}

		switch n.Kind {
		case KindError, KindComplete:
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

	switch res.Kind {
	case KindNext:
		return res.Value, nil
	case KindError:
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
//
// Like any other Blocking methods, this method waits for every goroutine
// started during subscription to complete before returning.
// To have this work properly, Observables must use [WaitGroupFromContext]
// to obtain a WaitGroup and use [WaitGroup.Go] rather than built-in go
// statements to start new goroutines during subscription, especially when
// they need to subscribe to other Observables in a goroutine; otherwise,
// runtime panicking might happen randomly (WaitGroup misuse).
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
//
// Like any other Blocking methods, this method waits for every goroutine
// started during subscription to complete before returning.
// To have this work properly, Observables must use [WaitGroupFromContext]
// to obtain a WaitGroup and use [WaitGroup.Go] rather than built-in go
// statements to start new goroutines during subscription, especially when
// they need to subscribe to other Observables in a goroutine; otherwise,
// runtime panicking might happen randomly (WaitGroup misuse).
func (obs Observable[T]) BlockingSingle(ctx context.Context) (v T, err error) {
	var wg WaitGroup

	child, cancel := context.WithCancel(WithWaitGroup(ctx, &wg))

	res := Error[T](ErrEmpty)

	var noop bool

	wg.Add(1)

	obs.Subscribe(child, func(n Notification[T]) {
		if noop {
			return
		}

		if n.Kind == KindNext && res.Kind == KindNext {
			res = Error[T](ErrNotSingle)
			noop = true

			cancel()
			wg.Done()

			return
		}

		switch n.Kind {
		case KindNext, KindError:
			res = n
		}

		switch n.Kind {
		case KindError, KindComplete:
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

	switch res.Kind {
	case KindNext:
		return res.Value, nil
	case KindError:
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
//
// Like any other Blocking methods, this method waits for every goroutine
// started during subscription to complete before returning.
// To have this work properly, Observables must use [WaitGroupFromContext]
// to obtain a WaitGroup and use [WaitGroup.Go] rather than built-in go
// statements to start new goroutines during subscription, especially when
// they need to subscribe to other Observables in a goroutine; otherwise,
// runtime panicking might happen randomly (WaitGroup misuse).
func (obs Observable[T]) BlockingSubscribe(ctx context.Context, sink Observer[T]) error {
	var wg WaitGroup

	child := WithWaitGroup(ctx, &wg)

	var res Notification[T]

	wg.Add(1)

	obs.Subscribe(child, func(n Notification[T]) {
		res = n

		sink(n)

		switch n.Kind {
		case KindError, KindComplete:
			wg.Done()
		}
	})

	wg.Wait()

	switch res.Kind {
	case KindError:
		return res.Error
	case KindComplete:
		return nil
	default:
		panic("unreachable")
	}
}
