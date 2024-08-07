package rx

// Connect multicasts the source Observable within a function where multiple
// subscriptions can share the same source.
func Connect[T, R any](selector func(source Observable[T]) Observable[R]) ConnectOperator[T, R] {
	return ConnectOperator[T, R]{
		ts: connectConfig[T, R]{
			connector: Multicast[T],
			selector:  selector,
		},
	}
}

type connectConfig[T, R any] struct {
	connector func() Subject[T]
	selector  func(Observable[T]) Observable[R]
}

// ConnectOperator is an [Operator] type for [Connect].
type ConnectOperator[T, R any] struct {
	ts connectConfig[T, R]
}

// WithConnector sets Connector option to a given value.
func (op ConnectOperator[T, R]) WithConnector(connector func() Subject[T]) ConnectOperator[T, R] {
	op.ts.connector = connector
	return op
}

// Apply implements the Operator interface.
func (op ConnectOperator[T, R]) Apply(source Observable[T]) Observable[R] {
	return connectObservable[T, R]{source, op.ts}.Subscribe
}

type connectObservable[T, R any] struct {
	source Observable[T]
	connectConfig[T, R]
}

func (ob connectObservable[T, R]) Subscribe(c Context, o Observer[R]) {
	c, cancel := c.WithCancel()
	o = o.DoOnTermination(cancel)
	oops := func() { o.Error(ErrOops) }
	subject := Try01(ob.connector, oops)
	Try11(ob.selector, subject.Observable, oops).Subscribe(c, o)
	ob.source.Subscribe(c, subject.Observer)
}
