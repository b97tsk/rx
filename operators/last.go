package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Last creates an Observable that emits only the last value emitted by
// the source Observable, if the source turns out to be empty, throws
// rx.ErrEmpty.
func Last() rx.Operator {
	return last
}

func last(source rx.Observable) rx.Observable {
	return func(ctx context.Context, sink rx.Observer) {
		var last struct {
			Value    interface{}
			HasValue bool
		}
		source.Subscribe(ctx, func(t rx.Notification) {
			switch {
			case t.HasValue:
				last.Value = t.Value
				last.HasValue = true
			case t.HasError:
				sink(t)
			default:
				if last.HasValue {
					sink.Next(last.Value)
					sink.Complete()
				} else {
					sink.Error(rx.ErrEmpty)
				}
			}
		})
	}
}

// LastOrDefault creates an Observable that emits only the last value emitted
// by the source Observable, if the source turns out to be empty, emits a
// specified default value.
func LastOrDefault(def interface{}) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			var last struct {
				Value    interface{}
				HasValue bool
			}
			source.Subscribe(ctx, func(t rx.Notification) {
				switch {
				case t.HasValue:
					last.Value = t.Value
					last.HasValue = true
				case t.HasError:
					sink(t)
				default:
					if last.HasValue {
						sink.Next(last.Value)
					} else {
						sink.Next(def)
					}
					sink(t)
				}
			})
		}
	}
}
