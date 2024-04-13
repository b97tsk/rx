package rx

// BlockingFirst subscribes to ob, returning the first emitted value.
// If ob emits no values, it returns the zero value of T and ErrEmpty;
// if ob emits a notification of error, it returns the zero value of T and
// the error.
//
// The cancellation of parent will cause BlockingFirst to immediately return
// the zero value of T and parent.Cause().
func (ob Observable[T]) BlockingFirst(parent Context) (v T, err error) {
	res := Error[T](ErrEmpty)
	c, cancel := parent.WithCancel()

	var noop bool

	ob.Subscribe(c, func(n Notification[T]) {
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
		}
	})

	<-c.Done()

	select {
	default:
	case <-parent.Done():
		return v, parent.Cause()
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

// BlockingFirstOrElse subscribes to ob, returning the first emitted value or
// def if ob emits no values or emits a notification of error.
//
// The cancellation of parent will cause BlockingFirstOrElse to immediately
// return def.
func (ob Observable[T]) BlockingFirstOrElse(parent Context, def T) T {
	v, err := ob.BlockingFirst(parent)
	if err != nil {
		return def
	}

	return v
}

// BlockingLast subscribes to ob, returning the last emitted value.
// If ob emits no values, it returns the zero value of T and ErrEmpty;
// if ob emits a notification of error, it returns the zero value of T and
// the error.
//
// The cancellation of parent will cause BlockingLast to immediately return
// the zero value of T and parent.Cause().
func (ob Observable[T]) BlockingLast(parent Context) (v T, err error) {
	res := Error[T](ErrEmpty)
	c, cancel := parent.WithCancel()

	ob.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext, KindError:
			res = n
		}

		switch n.Kind {
		case KindError, KindComplete:
			cancel()
		}
	})

	<-c.Done()

	select {
	default:
	case <-parent.Done():
		return v, parent.Cause()
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

// BlockingLastOrElse subscribes to ob, returning the last emitted value or
// def if ob emits no values or emits a notification of error.
//
// The cancellation of parent will cause BlockingLastOrElse to immediately
// return def.
func (ob Observable[T]) BlockingLastOrElse(parent Context, def T) T {
	v, err := ob.BlockingLast(parent)
	if err != nil {
		return def
	}

	return v
}

// BlockingSingle subscribes to ob, returning the single emitted value.
// If ob emits more than one value or no values, it returns the zero value of
// T and ErrNotSingle or ErrEmpty respectively; if ob emits a notification of
// error, it returns the zero value of T and the error.
//
// The cancellation of parent will cause BlockingSingle to immediately return
// the zero value of T and parent.Cause().
func (ob Observable[T]) BlockingSingle(parent Context) (v T, err error) {
	res := Error[T](ErrEmpty)
	c, cancel := parent.WithCancel()

	var noop bool

	ob.Subscribe(c, func(n Notification[T]) {
		if noop {
			return
		}

		if n.Kind == KindNext && res.Kind == KindNext {
			res = Error[T](ErrNotSingle)
			noop = true
			cancel()
			return
		}

		switch n.Kind {
		case KindNext, KindError:
			res = n
		}

		switch n.Kind {
		case KindError, KindComplete:
			cancel()
		}
	})

	<-c.Done()

	select {
	default:
	case <-parent.Done():
		return v, parent.Cause()
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

// BlockingSubscribe subscribes to ob and waits for it to complete.
// If ob completes without an error, BlockingSubscribe returns nil;
// otherwise, it returns the emitted error.
//
// The cancellation of parent will cause BlockingSubscribe to immediately
// return parent.Cause().
func (ob Observable[T]) BlockingSubscribe(parent Context, o Observer[T]) error {
	var res Notification[T]

	c, cancel := parent.WithCancel()

	ob.Subscribe(c, func(n Notification[T]) {
		res = n
		o.Emit(n)
		switch n.Kind {
		case KindError, KindComplete:
			cancel()
		}
	})

	<-c.Done()

	select {
	default:
	case <-parent.Done():
		return parent.Cause()
	}

	switch res.Kind {
	case KindError:
		return res.Error
	case KindComplete:
		return nil
	default:
		panic("unreachable")
	}
}
