package rx

// Connect multicasts the source Observable within a function where multiple
// subscriptions can share the same source.
func Connect[T, R any](selector func(source Observable[T]) Observable[R]) ConnectOperator[T, R] {
	return ConnectOperator[T, R]{
		ts: connectConfig[T, R]{
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
	ts connectConfig[T, R]
}

// WithConnector sets Connector option to a given value.
func (op ConnectOperator[T, R]) WithConnector(connector func() Subject[T]) ConnectOperator[T, R] {
	op.ts.Connector = connector
	return op
}

// Apply implements the Operator interface.
func (op ConnectOperator[T, R]) Apply(source Observable[T]) Observable[R] {
	return connectObservable[T, R]{source, op.ts}.Subscribe
}

type connectObservable[T, R any] struct {
	Source Observable[T]
	connectConfig[T, R]
}

func (ob connectObservable[T, R]) Subscribe(c Context, o Observer[R]) {
	c, cancel := c.WithCancel()
	o = o.DoOnTermination(cancel)
	oops := func() { o.Error(ErrOops) }
	subject := Try01(ob.Connector, oops)
	Try11(ob.Selector, subject.Observable, oops).Subscribe(c, o)
	ob.Source.Subscribe(c, subject.Observer)
}
