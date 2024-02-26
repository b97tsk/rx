package rx

// Every emits a boolean to indicate whether every value of the source
// Observable satisfies a given condition.
func Every[T any](cond func(v T) bool) Operator[T, bool] {
	if cond == nil {
		panic("cond == nil")
	}

	return every(cond)
}

func every[T any](cond func(v T) bool) Operator[T, bool] {
	return NewOperator(
		func(source Observable[T]) Observable[bool] {
			return everyObservable[T]{source, cond}.Subscribe
		},
	)
}

type everyObservable[T any] struct {
	Source    Observable[T]
	Condition func(T) bool
}

func (obs everyObservable[T]) Subscribe(c Context, sink Observer[bool]) {
	c, cancel := c.WithCancel()
	sink = sink.OnLastNotification(cancel)

	var noop bool

	obs.Source.Subscribe(c, func(n Notification[T]) {
		if noop {
			return
		}

		switch n.Kind {
		case KindNext:
			if !obs.Condition(n.Value) {
				noop = true
				sink.Next(false)
				sink.Complete()
			}
		case KindError:
			sink.Error(n.Error)
		case KindComplete:
			sink.Next(true)
			sink.Complete()
		}
	})
}
