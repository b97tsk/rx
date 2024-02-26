package rx

// Find emits only the first value emitted by the source Observable that
// meets some condition.
func Find[T any](cond func(v T) bool) Operator[T, T] {
	if cond == nil {
		panic("cond == nil")
	}

	return find(cond)
}

func find[T any](cond func(v T) bool) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return findObservable[T]{source, cond}.Subscribe
		},
	)
}

type findObservable[T any] struct {
	Source    Observable[T]
	Condition func(T) bool
}

func (obs findObservable[T]) Subscribe(c Context, sink Observer[T]) {
	c, cancel := c.WithCancel()
	sink = sink.OnLastNotification(cancel)

	var noop bool

	obs.Source.Subscribe(c, func(n Notification[T]) {
		if noop {
			return
		}

		switch n.Kind {
		case KindNext:
			if obs.Condition(n.Value) {
				noop = true
				sink(n)
				sink.Complete()
			}
		case KindError, KindComplete:
			sink(n)
		}
	})
}
