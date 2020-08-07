package rx

import (
	"context"
)

// BlockingFirst subscribes to the source Observable, returns the first item
// emitted by the source; if the source emits no items, it returns nil plus
// ErrEmpty; if the source errors, it returns nil plus the error.
func (obs Observable) BlockingFirst(ctx context.Context) (interface{}, error) {
	var (
		result   Notification
		observer Observer
	)
	childCtx, childCancel := context.WithCancel(ctx)
	observer = func(t Notification) {
		switch {
		case t.HasValue || t.HasError:
			result = t
		default:
			result = Notification{Error: ErrEmpty, HasError: true}
		}
		observer = Noop
		childCancel()
	}
	obs.Subscribe(childCtx, observer.Sink)
	<-childCtx.Done()
	switch {
	case ctx.Err() != nil:
		return nil, ctx.Err()
	case result.HasValue:
		return result.Value, nil
	case result.HasError:
		return nil, result.Error
	default:
		return nil, childCtx.Err()
	}
}

// BlockingFirstOrDefault subscribes to the source Observable, returns the
// first item emitted by the source, or returns def if the source emits no
// items or emits an error.
func (obs Observable) BlockingFirstOrDefault(ctx context.Context, def interface{}) interface{} {
	val, err := obs.BlockingFirst(ctx)
	if err != nil {
		return def
	}
	return val
}

// BlockingLast subscribes to the source Observable, returns the last item
// emitted by the source; if the source emits no items, it returns nil plus
// ErrEmpty; if the source errors, it returns nil plus the error.
func (obs Observable) BlockingLast(ctx context.Context) (interface{}, error) {
	var result Notification
	childCtx, childCancel := context.WithCancel(ctx)
	obs.Subscribe(childCtx, func(t Notification) {
		switch {
		case t.HasValue || t.HasError:
			result = t
		default:
			if !result.HasValue {
				result = Notification{Error: ErrEmpty, HasError: true}
			}
		}
		if !t.HasValue {
			childCancel()
		}
	})
	<-childCtx.Done()
	switch {
	case ctx.Err() != nil:
		return nil, ctx.Err()
	case result.HasValue:
		return result.Value, nil
	case result.HasError:
		return nil, result.Error
	default:
		return nil, childCtx.Err()
	}
}

// BlockingLastOrDefault subscribes to the source Observable, returns the last
// item emitted by the source, or returns def if the source emits no items or
// emits an error.
func (obs Observable) BlockingLastOrDefault(ctx context.Context, def interface{}) interface{} {
	val, err := obs.BlockingLast(ctx)
	if err != nil {
		return def
	}
	return val
}

// BlockingSingle subscribes to the source Observable, returns the single item
// emitted by the source; if the source emits more than one item or no items,
// it returns nil plus ErrNotSingle or ErrEmpty respectively; if the source
// errors, it returns nil plus the error.
func (obs Observable) BlockingSingle(ctx context.Context) (interface{}, error) {
	var (
		result   Notification
		observer Observer
	)
	childCtx, childCancel := context.WithCancel(ctx)
	observer = func(t Notification) {
		switch {
		case t.HasValue:
			if result.HasValue {
				result = Notification{Error: ErrNotSingle, HasError: true}
				observer = Noop
				childCancel()
			} else {
				result = t
			}
		case t.HasError:
			result = t
		default:
			if !result.HasValue {
				result = Notification{Error: ErrEmpty, HasError: true}
			}
		}
		if !t.HasValue {
			childCancel()
		}
	}
	obs.Subscribe(childCtx, observer.Sink)
	<-childCtx.Done()
	switch {
	case ctx.Err() != nil:
		return nil, ctx.Err()
	case result.HasValue:
		return result.Value, nil
	case result.HasError:
		return nil, result.Error
	default:
		return nil, childCtx.Err()
	}
}

// BlockingSubscribe subscribes to the source Observable, returns only when
// the source completes or errors; if the source completes, it returns nil;
// if the source errors, it returns the error.
//
// Note that sink may be called even after BlockingSubscribe has returned.
// This only happens if ctx was cancelled during the subscription. A possible
// workaround is:
//
//	obs.BlockingSubscribe(ctx, func(t Notification) {
//		if ctx.Err() != nil {
//			return
//		}
//		/* rest of the code */
//	})
//
func (obs Observable) BlockingSubscribe(ctx context.Context, sink Observer) error {
	var result Notification
	childCtx, childCancel := context.WithCancel(ctx)
	obs.Subscribe(childCtx, func(t Notification) {
		if !t.HasValue {
			defer childCancel()
		}
		result = t
		sink(t)
	})
	<-childCtx.Done()
	switch {
	case ctx.Err() != nil:
		return ctx.Err()
	case result.HasValue:
		return childCtx.Err()
	case result.HasError:
		return result.Error
	default:
		return nil
	}
}
