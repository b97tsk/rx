package rx

import "sync"

// BlockingFirst subscribes to obs, returning the first emitted value.
// If obs emits no values, it returns the zero value of T and ErrEmpty;
// if obs emits an error notification, it returns the zero value of T and
// the error.
//
// obs must honor the cancellation of c; otherwise, BlockingFirst might
// continue to block even after c has been canceled.
//
// Like any other Blocking methods, this method sets c.WaitGroup to
// a new [sync.WaitGroup] and calls its Wait method after subscribing to obs.
// This method returns only when the WaitGroup counter is zero.
func (obs Observable[T]) BlockingFirst(c Context) (v T, err error) {
	res := Error[T](ErrEmpty)

	var wg sync.WaitGroup

	c, cancel := c.WithWaitGroup(&wg).WithCancel()

	var noop bool

	wg.Add(1)

	obs.Subscribe(c, func(n Notification[T]) {
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

	switch res.Kind {
	case KindNext:
		return res.Value, nil
	case KindError:
		return v, res.Error
	default:
		panic("unreachable")
	}
}

// BlockingFirstOrElse subscribes to obs, returning the first emitted value or
// def if obs emits no values or an error notification.
//
// obs must honor the cancellation of c; otherwise, BlockingFirstOrElse might
// continue to block even after c has been canceled.
//
// Like any other Blocking methods, this method sets c.WaitGroup to
// a new [sync.WaitGroup] and calls its Wait method after subscribing to obs.
// This method returns only when the WaitGroup counter is zero.
func (obs Observable[T]) BlockingFirstOrElse(c Context, def T) T {
	v, err := obs.BlockingFirst(c)
	if err != nil {
		return def
	}

	return v
}

// BlockingLast subscribes to obs, returning the last emitted value.
// If obs emits no values, it returns the zero value of T and ErrEmpty;
// if obs emits an error notification, it returns the zero value of T and
// the error.
//
// obs must honor the cancellation of c; otherwise, BlockingLast might
// continue to block even after c has been canceled.
//
// Like any other Blocking methods, this method sets c.WaitGroup to
// a new [sync.WaitGroup] and calls its Wait method after subscribing to obs.
// This method returns only when the WaitGroup counter is zero.
func (obs Observable[T]) BlockingLast(c Context) (v T, err error) {
	res := Error[T](ErrEmpty)

	var wg sync.WaitGroup

	c, cancel := c.WithWaitGroup(&wg).WithCancel()

	wg.Add(1)

	obs.Subscribe(c, func(n Notification[T]) {
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

	switch res.Kind {
	case KindNext:
		return res.Value, nil
	case KindError:
		return v, res.Error
	default:
		panic("unreachable")
	}
}

// BlockingLastOrElse subscribes to obs, returning the last emitted value or
// def if obs emits no values or an error notification.
//
// obs must honor the cancellation of c; otherwise, BlockingLastOrElse might
// continue to block even after c has been canceled.
//
// Like any other Blocking methods, this method sets c.WaitGroup to
// a new [sync.WaitGroup] and calls its Wait method after subscribing to obs.
// This method returns only when the WaitGroup counter is zero.
func (obs Observable[T]) BlockingLastOrElse(c Context, def T) T {
	v, err := obs.BlockingLast(c)
	if err != nil {
		return def
	}

	return v
}

// BlockingSingle subscribes to obs, returning the single emitted value.
// If obs emits more than one value or no values, it returns the zero value of
// T and ErrNotSingle or ErrEmpty respectively; if obs emits a notification of
// error, it returns the zero value of T and the error.
//
// obs must honor the cancellation of c; otherwise, BlockingSingle might
// continue to block even after c has been canceled.
//
// Like any other Blocking methods, this method sets c.WaitGroup to
// a new [sync.WaitGroup] and calls its Wait method after subscribing to obs.
// This method returns only when the WaitGroup counter is zero.
func (obs Observable[T]) BlockingSingle(c Context) (v T, err error) {
	res := Error[T](ErrEmpty)

	var wg sync.WaitGroup

	c, cancel := c.WithWaitGroup(&wg).WithCancel()

	var noop bool

	wg.Add(1)

	obs.Subscribe(c, func(n Notification[T]) {
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

	switch res.Kind {
	case KindNext:
		return res.Value, nil
	case KindError:
		return v, res.Error
	default:
		panic("unreachable")
	}
}

// BlockingSubscribe subscribes to obs and waits for it to complete.
// If obs completes without an error, BlockingSubscribe returns nil;
// otherwise, it returns the emitted error.
//
// obs must honor the cancellation of c; otherwise, BlockingSubscribe might
// continue to block even after c has been canceled.
//
// Like any other Blocking methods, this method sets c.WaitGroup to
// a new [sync.WaitGroup] and calls its Wait method after subscribing to obs.
// This method returns only when the WaitGroup counter is zero.
func (obs Observable[T]) BlockingSubscribe(c Context, sink Observer[T]) error {
	var res Notification[T]

	var wg sync.WaitGroup

	c = c.WithWaitGroup(&wg)

	wg.Add(1)

	obs.Subscribe(c, func(n Notification[T]) {
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
