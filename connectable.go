package rx

import (
	"context"
	"sync"
	"time"
)

// A ConnectableObservable is an Observable that only subscribes to the source
// Observable by calling its Connect method. Calling its Subscribe method will
// not subscribe the source, instead, it subscribes to a local Subject, which
// means that it can be called many times with different Observers.
type ConnectableObservable struct {
	*connectableObservable
}

type connectableObservable struct {
	Observable
	mux            sync.Mutex
	source         Observable
	subjectFactory func() Subject
	subject        Subject
	connection     context.Context
	disconnect     context.CancelFunc
}

func newConnectableObservable(source Observable, subjectFactory func() Subject) ConnectableObservable {
	obs := &connectableObservable{
		source:         source,
		subjectFactory: subjectFactory,
	}
	obs.Observable = Observable(
		func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
			return obs.getSubject().Subscribe(ctx, sink)
		},
	)
	return ConnectableObservable{obs}
}

// Exists reports if this ConnectableObservable is ready to use.
func (obs ConnectableObservable) Exists() bool {
	return obs.connectableObservable != nil
}

func (obs *connectableObservable) getSubject() Subject {
	obs.mux.Lock()
	defer obs.mux.Unlock()
	return obs.getSubjectLocked()
}

func (obs *connectableObservable) getSubjectLocked() Subject {
	if obs.subject.Observer == nil {
		obs.subject = obs.subjectFactory()
	}
	return obs.subject
}

func (obs *connectableObservable) connect(ctx context.Context) (context.Context, context.CancelFunc) {
	obs.mux.Lock()
	defer obs.mux.Unlock()

	connection := obs.connection

	if connection == nil {
		type X struct{}
		cx := make(chan X, 1)
		cx <- X{}

		sink := obs.getSubjectLocked().Observer

		ctx, cancel := obs.source.Subscribe(ctx, func(t Notification) {
			if t.HasValue {
				sink(t)
				return
			}

			x, cxLocked := <-cx
			if !cxLocked {
				obs.mux.Lock()
			}

			if connection == obs.connection {
				obs.subject = Subject{}
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
			obs.subject = Subject{}
			obs.connection = nil
			obs.disconnect = nil
		}
		obs.mux.Unlock()
	}
}

// Connect invokes an execution of an ConnectableObservable.
func (obs ConnectableObservable) Connect(ctx context.Context) (context.Context, context.CancelFunc) {
	return obs.connect(ctx)
}

// Multicast returns a ConnectableObservable, which is a variety of Observable
// that waits until its Connect method is called before it begins emitting
// items to those Observers that have subscribed to it.
func (obs Observable) Multicast(subjectFactory func() Subject) ConnectableObservable {
	return newConnectableObservable(obs, subjectFactory)
}

// Publish is like Multicast, but it uses only one subject.
func (obs Observable) Publish() ConnectableObservable {
	subject := NewSubject()
	return obs.Multicast(func() Subject { return subject })
}

// PublishBehavior is like Publish, but it uses a BehaviorSubject instead.
func (obs Observable) PublishBehavior(val interface{}) ConnectableObservable {
	subject := NewBehaviorSubject(val)
	return obs.Multicast(func() Subject { return subject.Subject })
}

// PublishReplay is like Publish, but it uses a ReplaySubject instead.
func (obs Observable) PublishReplay(bufferSize int, windowTime time.Duration) ConnectableObservable {
	subject := NewReplaySubject(bufferSize, windowTime)
	return obs.Multicast(func() Subject { return subject.Subject })
}
