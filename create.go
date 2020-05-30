package rx

import (
	"context"
	"errors"
	"sync/atomic"
)

var errComplete = errors.New("complete")

// Complete is a special error that represents a complete subscription.
var Complete = errComplete

// Create creates a new Observable, that will execute the specified function
// when an Observer subscribes to it.
//
// It's the caller's responsibility to follow the Observable Contract that
// no more emissions pass to the sink Observer after an ERROR or COMPLETE
// emission has passed to it. Violations of this contract result in undefined
// behaviors.
func Create(subscribe func(context.Context, Observer)) Observable {
	return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
		ctx, cancel := context.WithCancel(ctx)
		k := &kontext{Context: ctx}
		subscribe(k, func(t Notification) {
			sink(t)
			switch {
			case t.HasValue:
			case t.HasError:
				k.err.Store(t.Error)
				cancel()
			default:
				k.err.Store(Complete)
				cancel()
			}
		})
		return k, cancel
	}
}

type kontext struct {
	context.Context
	err atomic.Value
}

func (k *kontext) Err() error {
	if err := k.err.Load(); err != nil {
		return err.(error)
	}
	return k.Context.Err()
}
