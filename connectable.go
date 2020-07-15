package rx

import (
	"context"
	"sync"
)

// A ConnectableObservable is an Observable that only subscribes to the source
// Observable by calling its Connect method. Calling its Subscribe method will
// not subscribe the source, instead, it subscribes to a local Double, which
// means that it can be called many times with different Observers.
type ConnectableObservable struct {
	Observable

	mux           sync.Mutex
	source        Observable
	doubleFactory func() Double
	double        Double
	connection    context.Context
	disconnect    context.CancelFunc
}

func newConnectableObservable(source Observable, doubleFactory func() Double) *ConnectableObservable {
	obs := &ConnectableObservable{
		source:        source,
		doubleFactory: doubleFactory,
	}
	obs.Observable = Observable(
		func(ctx context.Context, sink Observer) {
			obs.getDouble().Subscribe(ctx, sink)
		},
	)
	return obs
}

func (obs *ConnectableObservable) getDouble() Double {
	obs.mux.Lock()
	defer obs.mux.Unlock()
	return obs.getDoubleLocked()
}

func (obs *ConnectableObservable) getDoubleLocked() Double {
	if obs.double.Observable == nil {
		obs.double = obs.doubleFactory()
	}
	return obs.double
}

// Connect invokes an execution of an ConnectableObservable.
func (obs *ConnectableObservable) Connect(ctx context.Context) (context.Context, context.CancelFunc) {
	obs.mux.Lock()
	defer obs.mux.Unlock()

	connection := obs.connection

	if connection == nil {
		type X struct{}
		cx := make(chan X, 1)
		cx <- X{}

		sink := obs.getDoubleLocked().Observer
		ctx, cancel := context.WithCancel(ctx)
		obs.source.Subscribe(ctx, func(t Notification) {
			if t.HasValue {
				sink(t)
				return
			}

			cancel()

			x, cxLocked := <-cx
			if !cxLocked {
				obs.mux.Lock()
			}

			if connection == obs.connection {
				obs.double = Double{}
				obs.connection = nil
				obs.disconnect = nil
			}

			if cxLocked {
				cx <- x
			} else {
				obs.mux.Unlock()
			}

			sink(t)
		})

		<-cx
		close(cx)

		if ctx.Err() != nil {
			return ctx, cancel
		}

		connection = ctx
		obs.connection = ctx
		obs.disconnect = cancel
	}

	return connection, func() {
		obs.mux.Lock()
		if connection == obs.connection {
			obs.disconnect()
			obs.double = Double{}
			obs.connection = nil
			obs.disconnect = nil
		}
		obs.mux.Unlock()
	}
}

// Multicast returns a ConnectableObservable, which is a variety of Observable
// that waits until its Connect method is called before it begins emitting
// items to those Observers that have subscribed to it.
func (obs Observable) Multicast(doubleFactory func() Double) *ConnectableObservable {
	return newConnectableObservable(obs, doubleFactory)
}

// Publish is like Multicast, but it uses only one Double.
func (obs Observable) Publish(d Double) *ConnectableObservable {
	return obs.Multicast(func() Double { return d })
}
