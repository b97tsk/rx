package rx

import (
	"context"
)

// Connect multicasts the source Observable within a function where multiple
// subscriptions can share the same source.
func Connect[T, R any](selector func(source Observable[T]) Observable[R]) ConnectOperator[T, R] {
	if selector == nil {
		panic("selector == nil")
	}

	return ConnectOperator[T, R]{
		opts: connectConfig[T, R]{
			Connector: Multicast[T],
			Selector:  selector,
		},
	}
}

type connectConfig[T, R any] struct {
	Connector func() Subject[T]
	Selector  func(Observable[T]) Observable[R]
}

// ConnectOperator is an [Operator] type for [Connect].
type ConnectOperator[T, R any] struct {
	opts connectConfig[T, R]
}

// WithConnector sets Connector option to a given value.
func (op ConnectOperator[T, R]) WithConnector(connector func() Subject[T]) ConnectOperator[T, R] {
	if connector == nil {
		panic("connector == nil")
	}

	op.opts.Connector = connector

	return op
}

// Apply implements the Operator interface.
func (op ConnectOperator[T, R]) Apply(source Observable[T]) Observable[R] {
	return connectObservable[T, R]{source, op.opts}.Subscribe
}

type connectObservable[T, R any] struct {
	Source Observable[T]
	connectConfig[T, R]
}

func (obs connectObservable[T, R]) Subscribe(ctx context.Context, sink Observer[R]) {
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.OnLastNotification(cancel)
	subject := obs.Connector()
	obs.Selector(subject.Observable).Subscribe(ctx, sink)
	obs.Source.Subscribe(ctx, subject.Observer)
}
