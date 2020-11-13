package rx

import (
	"context"
)

// BlockingFirst subscribes to the source Observable, returns the first item
// emitted by the source; if the source emits no items, it returns nil plus
// ErrEmpty; if the source emits an error, it returns nil plus this error.
//
// If ctx was cancelled during the subscription, BlockingFirst immediately
// returns nil plus ctx.Err().
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
	default: // Unreachable path.
		return nil, childCtx.Err()
	}
}

// BlockingFirstOrDefault subscribes to the source Observable, returns the
// first item emitted by the source, or returns def if the source emits no
// items or an error.
//
// If ctx was cancelled during the subscription, BlockingFirstOrDefault
// immediately returns def.
func (obs Observable) BlockingFirstOrDefault(ctx context.Context, def interface{}) interface{} {
	val, err := obs.BlockingFirst(ctx)
	if err != nil {
		return def
	}
	return val
}

// BlockingLast subscribes to the source Observable, returns the last item
// emitted by the source; if the source emits no items, it returns nil plus
// ErrEmpty; if the source emits an error, it returns nil plus this error.
//
// If ctx was cancelled during the subscription, BlockingLast immediately
// returns nil plus ctx.Err().
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
	default: // Unreachable path.
		return nil, childCtx.Err()
	}
}

// BlockingLastOrDefault subscribes to the source Observable, returns the last
// item emitted by the source, or returns def if the source emits no items or
// an error.
//
// If ctx was cancelled during the subscription, BlockingLastOrDefault
// immediately returns def.
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
// emits an error, it returns nil plus this error.
//
// If ctx was cancelled during the subscription, BlockingSingle immediately
// returns nil plus ctx.Err().
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
	default: // Unreachable path.
		return nil, childCtx.Err()
	}
}

// BlockingSubscribe subscribes to the source Observable, returns only when
// the source completes or emits an error; if the source completes, it returns
// nil; if the source emits an error, it returns this error.
//
// If ctx was cancelled during the subscription, BlockingSubscribe immediately
// returns ctx.Err() and sink may still be called after that.
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
	case result.HasValue: // Always false.
		return childCtx.Err()
	case result.HasError:
		return result.Error
	default:
		return nil
	}
}
