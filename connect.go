package rx

// Connect multicasts the source Observable within a function where multiple
// subscriptions can share the same source.
func Connect[T, R any](selector func(source Observable[T]) Observable[R]) ConnectOperator[T, R] {
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

func (obs connectObservable[T, R]) Subscribe(c Context, sink Observer[R]) {
	c, cancel := c.WithCancel()
	sink = sink.OnLastNotification(cancel)
	oops := func() { sink.Error(ErrOops) }
	subject := Try01(obs.Connector, oops)
	Try11(obs.Selector, subject.Observable, oops).Subscribe(c, sink)
	obs.Source.Subscribe(c, subject.Observer)
}
