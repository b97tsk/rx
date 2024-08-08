package rx

import "sync"

// BlockingFirst subscribes to ob, returning the first value notification.
// If ob completes without emitting any value, BlockingFirst returns
// an [Error] notification of [ErrEmpty];
// if ob emits no values but a notification of [Error] or [Stop],
// BlockingFirst returns that notification.
func (ob Observable[T]) BlockingFirst(parent Context) Notification[T] {
	c, cancel := parent.WithCancel()

	var x struct {
		wg   sync.WaitGroup
		res  Notification[T]
		noop bool
	}

	x.wg.Add(1)

	ob.Subscribe(c, func(n Notification[T]) {
		if x.noop {
			return
		}

		switch n.Kind {
		case KindNext:
			x.noop = true
			x.res = n
			cancel()
			x.wg.Done()
		case KindComplete:
			x.res = Error[T](ErrEmpty)
			cancel()
			x.wg.Done()
		case KindError, KindStop:
			x.res = n
			cancel()
			x.wg.Done()
		}
	})

	x.wg.Wait()

	return x.res
}

// BlockingLast subscribes to ob, returning the last value notification.
// If ob completes without emitting any value, BlockingLast returns
// an [Error] notification of [ErrEmpty];
// if ob emits a notification of [Error] or [Stop], BlockingLast returns
// that notification.
func (ob Observable[T]) BlockingLast(c Context) Notification[T] {
	var x struct {
		wg  sync.WaitGroup
		res Notification[T]
	}

	x.wg.Add(1)

	ob.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			x.res = n
		case KindComplete:
			if x.res.Kind == 0 {
				x.res = Error[T](ErrEmpty)
			}
			x.wg.Done()
		case KindError, KindStop:
			x.res = n
			x.wg.Done()
		}
	})

	x.wg.Wait()

	return x.res
}

// BlockingSingle subscribes to ob, returning the single value notification.
// If ob emits more than one value, BlockingSingle returns an [Error]
// notification of [ErrNotSingle];
// if ob completes without emitting any value, BlockingSingle returns
// an [Error] notification of [ErrEmpty];
// if ob emits a notification of [Error] or [Stop], BlockingSingle returns
// that notification.
func (ob Observable[T]) BlockingSingle(parent Context) Notification[T] {
	c, cancel := parent.WithCancel()

	var x struct {
		wg   sync.WaitGroup
		res  Notification[T]
		noop bool
	}

	x.wg.Add(1)

	ob.Subscribe(c, func(n Notification[T]) {
		if x.noop {
			return
		}

		if n.Kind == KindNext && x.res.Kind == KindNext {
			x.res = Error[T](ErrNotSingle)
			x.noop = true
			cancel()
			x.wg.Done()
			return
		}

		switch n.Kind {
		case KindNext:
			x.res = n
		case KindComplete:
			if x.res.Kind == 0 {
				x.res = Error[T](ErrEmpty)
			}
			cancel()
			x.wg.Done()
		case KindError, KindStop:
			x.res = n
			cancel()
			x.wg.Done()
		}
	})

	x.wg.Wait()

	return x.res
}

// BlockingSubscribe subscribes to ob and waits for it to complete,
// returning the last notification which must be one of [Complete], [Error]
// or [Stop].
func (ob Observable[T]) BlockingSubscribe(c Context, o Observer[T]) Notification[T] {
	var x struct {
		wg  sync.WaitGroup
		res Notification[T]
	}

	x.wg.Add(1)

	ob.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
		case KindComplete, KindError, KindStop:
			defer x.wg.Done()
		}
		x.res = n
		o.Emit(n)
	})

	x.wg.Wait()

	return x.res
}
