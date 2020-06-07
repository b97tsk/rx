package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type elementAtObservable struct {
	Source     rx.Observable
	Index      int
	Default    interface{}
	HasDefault bool
}

func (obs elementAtObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var observer rx.Observer

	index := obs.Index

	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			index--
			if index == -1 {
				observer = rx.Noop
				sink(t)
				sink.Complete()
			}
		case t.HasError:
			sink(t)
		default:
			if obs.HasDefault {
				sink.Next(obs.Default)
				sink.Complete()
			} else {
				sink.Error(rx.ErrOutOfRange)
			}
		}
	}

	obs.Source.Subscribe(ctx, observer.Sink)
}

// ElementAt creates an Observable that emits the single value at the specified
// index in a sequence of emissions from the source Observable, if the
// specified index is out of range, notifies rx.ErrOutOfRange.
func ElementAt(index int) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := elementAtObservable{
			Source: source,
			Index:  index,
		}
		return rx.Create(obs.Subscribe)
	}
}

// ElementAtOrDefault creates an Observable that emits the single value at the
// specified index in a sequence of emissions from the source Observable, if
// the specified index is out of range, emits the provided default value.
func ElementAtOrDefault(index int, def interface{}) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := elementAtObservable{
			Source:     source,
			Index:      index,
			Default:    def,
			HasDefault: true,
		}
		return rx.Create(obs.Subscribe)
	}
}
