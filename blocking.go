package rx

import (
	"context"
)

// BlockingFirst subscribes to the source Observable, returns the first item
// emitted by the source; if the source emits no items, it returns nil and
// ErrEmpty; if the source errors, it returns nil and the error.
func (obs Observable) BlockingFirst(ctx context.Context) (interface{}, error) {
	var (
		result   Notification
		observer Observer
	)
	ctx, cancel := context.WithCancel(ctx)
	observer = func(t Notification) {
		observer = Noop
		switch {
		case t.HasValue || t.HasError:
			result = t
		default:
			result = Notification{Error: ErrEmpty, HasError: true}
		}
		cancel()
	}
	ctx, _ = obs.Subscribe(ctx, observer.Notify)
	<-ctx.Done()
	switch {
	case result.HasValue:
		return result.Value, nil
	case result.HasError:
		return nil, result.Error
	default:
		return nil, ctx.Err()
	}
}

// BlockingLast subscribes to the source Observable, returns the last item
// emitted by the source; if the source emits no items, it returns nil and
// ErrEmpty; if the source errors, it returns nil and the error.
func (obs Observable) BlockingLast(ctx context.Context) (interface{}, error) {
	var result Notification
	ctx, _ = obs.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue || t.HasError:
			result = t
		default:
			if !result.HasValue {
				result = Notification{Error: ErrEmpty, HasError: true}
			}
		}
	})
	<-ctx.Done()
	switch {
	case result.HasValue:
		return result.Value, nil
	case result.HasError:
		return nil, result.Error
	default:
		return nil, ctx.Err()
	}
}

// BlockingSingle subscribes to the source Observable, returns the single item
// emitted by the source; if the source emits more than one item or no items,
// it returns nil and ErrNotSingle or ErrEmpty respectively; if the source
// errors, it returns nil and the error.
func (obs Observable) BlockingSingle(ctx context.Context) (interface{}, error) {
	var (
		result   Notification
		observer Observer
	)
	ctx, cancel := context.WithCancel(ctx)
	observer = func(t Notification) {
		switch {
		case t.HasValue:
			if result.HasValue {
				observer = Noop
				result = Notification{Error: ErrNotSingle, HasError: true}
				cancel()
			} else {
				result = t
			}
		case t.HasError:
			result = t
			cancel()
		default:
			if !result.HasValue {
				result = Notification{Error: ErrEmpty, HasError: true}
			}
			cancel()
		}
	}
	ctx, _ = obs.Subscribe(ctx, observer.Notify)
	<-ctx.Done()
	switch {
	case result.HasValue:
		return result.Value, nil
	case result.HasError:
		return nil, result.Error
	default:
		return nil, ctx.Err()
	}
}

// BlockingSubscribe subscribes to the source Observable, returns only when
// the source completes or errors; if the source completes, it returns nil;
// if the source errors, it returns the error.
func (obs Observable) BlockingSubscribe(ctx context.Context, sink Observer) error {
	var (
		result    Notification
		hasResult bool
	)
	ctx, _ = obs.Subscribe(ctx, func(t Notification) {
		result = t
		hasResult = true
		sink(t)
	})
	<-ctx.Done()
	switch {
	case !hasResult || result.HasValue:
		return ctx.Err()
	case result.HasError:
		return result.Error
	default:
		return nil
	}
}
