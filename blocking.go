package rx

import (
	"context"
)

// BlockingFirst subscribes to the source Observable, returns the first item
// emitted by the source; if the source emits no items, it returns nil and
// error ErrEmpty; if the source errors, it returns with the error.
func (o Observable) BlockingFirst(ctx context.Context) (value interface{}, err error) {
	ctx, cancel := context.WithCancel(ctx)

	var observer Observer
	observer = func(t Notification) {
		switch {
		case t.HasValue:
			observer = NopObserver
			value = t.Value
			cancel()
		case t.HasError:
			err = t.Value.(error)
			cancel()
		default:
			err = ErrEmpty
			cancel()
		}
	}
	o.Subscribe(ctx, observer.Notify)

	<-ctx.Done()
	return
}

// BlockingLast subscribes to the source Observable, returns the last item
// emitted by the source; if the source emits no items, it returns nil and
// error ErrEmpty; if the source errors, it returns with the error.
func (o Observable) BlockingLast(ctx context.Context) (value interface{}, err error) {
	var hasValue bool
	ctx, _ = o.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			value = t.Value
			hasValue = true
		case t.HasError:
			err = t.Value.(error)
		default:
			if !hasValue {
				err = ErrEmpty
			}
		}
	})
	<-ctx.Done()
	return
}

// BlockingSingle subscribes to the source Observable, returns the single item
// emitted by the source; if the source emits more than one item or no items,
// it returns with error ErrNotSingle or ErrEmpty respectively; if the source
// errors, it returns with the error.
func (o Observable) BlockingSingle(ctx context.Context) (value interface{}, err error) {
	ctx, cancel := context.WithCancel(ctx)

	var (
		hasValue bool
		observer Observer
	)

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			if hasValue {
				observer = NopObserver
				err = ErrNotSingle
				cancel()
			} else {
				value = t.Value
				hasValue = true
			}
		case t.HasError:
			err = t.Value.(error)
			cancel()
		default:
			if !hasValue {
				err = ErrEmpty
			}
			cancel()
		}
	}

	o.Subscribe(ctx, observer.Notify)
	<-ctx.Done()
	return
}

// BlockingSubscribe subscribes to the source Observable, returns only when
// the source completes or errors; if the source completes, it returns nil;
// if the source errors, it returns the error.
func (o Observable) BlockingSubscribe(ctx context.Context, sink Observer) (err error) {
	ctx, _ = o.Subscribe(ctx, func(t Notification) {
		if t.HasError {
			err = t.Value.(error)
		}
		sink(t)
	})
	<-ctx.Done()
	return
}
