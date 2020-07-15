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
	ctx, cancel := context.WithCancel(ctx)
	observer = func(t Notification) {
		switch {
		case t.HasValue || t.HasError:
			result = t
		default:
			result = Notification{Error: ErrEmpty, HasError: true}
		}
		observer = Noop
		cancel()
	}
	obs.Subscribe(ctx, observer.Sink)
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
// emitted by the source; if the source emits no items, it returns nil plus
// ErrEmpty; if the source errors, it returns nil plus the error.
func (obs Observable) BlockingLast(ctx context.Context) (interface{}, error) {
	var result Notification
	ctx, cancel := context.WithCancel(ctx)
	obs.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue || t.HasError:
			result = t
		default:
			if !result.HasValue {
				result = Notification{Error: ErrEmpty, HasError: true}
			}
		}
		if !t.HasValue {
			cancel()
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
// it returns nil plus ErrNotSingle or ErrEmpty respectively; if the source
// errors, it returns nil plus the error.
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
				result = Notification{Error: ErrNotSingle, HasError: true}
				observer = Noop
				cancel()
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
			cancel()
		}
	}
	obs.Subscribe(ctx, observer.Sink)
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
	ctx, cancel := context.WithCancel(ctx)
	obs.Subscribe(ctx, func(t Notification) {
		if !t.HasValue {
			defer cancel()
		}
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
