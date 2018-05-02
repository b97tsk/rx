package rx

import (
	"context"
)

// BlockingFirst subscribes the source Observable, returns the first item
// emitted by the source; if the source emits no items, it returns with an
// error ErrEmpty; if the source emits an error, it returns with that error.
func (o Observable) BlockingFirst(ctx context.Context) (value interface{}, err error) {
	ctx, cancel := context.WithCancel(ctx)

	mutable := MutableObserver{}

	mutable.Observer = ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			mutable.Observer = NopObserver
			value = t.Value
			cancel()
		case t.HasError:
			err = t.Value.(error)
			cancel()
		default:
			err = ErrEmpty
			cancel()
		}
	})

	o.Op.Call(ctx, &mutable)
	<-ctx.Done()
	return
}

// BlockingLast subscribes the source Observable, returns the last item emitted
// by the source; if the source emits no items, it returns with an error
// ErrEmpty; if the source emits an error, it returns with that error.
func (o Observable) BlockingLast(ctx context.Context) (value interface{}, err error) {
	ctx, _ = o.Op.Call(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			value = t.Value
		case t.HasError:
			err = t.Value.(error)
		default:
		}
	}))
	<-ctx.Done()
	return
}

// BlockingSingle subscribes the source Observable, returns the single item
// emitted by the source; if the source emits more than one item or no items,
// it returns with an error ErrNotSingle or ErrEmpty respectively; if the source
// emits an error, it returns with that error.
func (o Observable) BlockingSingle(ctx context.Context) (value interface{}, err error) {
	ctx, cancel := context.WithCancel(ctx)
	hasValue := false

	mutable := MutableObserver{}

	mutable.Observer = ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			if hasValue {
				mutable.Observer = NopObserver
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
	})

	o.Op.Call(ctx, &mutable)
	<-ctx.Done()
	return
}

// BlockingSubscribe subscribes the source Observable, returns only when the
// source completes or emits an error; if the source completes, it returns nil;
// if the source emits an error, it returns that error.
func (o Observable) BlockingSubscribe(ctx context.Context, ob Observer) (err error) {
	ctx, _ = o.Op.Call(ctx, ObserverFunc(func(t Notification) {
		if t.HasError {
			err = t.Value.(error)
		}
		t.Observe(ob)
	}))
	<-ctx.Done()
	return
}
