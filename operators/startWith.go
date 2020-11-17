package operators

import (
	"github.com/b97tsk/rx"
)

// StartWith creates an Observable that emits the items you specify as
// arguments before it begins to emit items emitted by the source Observable.
func StartWith(values ...interface{}) rx.Operator {
	if len(values) == 0 {
		return noop
	}

	return func(source rx.Observable) rx.Observable {
		return rx.Concat(rx.FromSlice(values), source)
	}
}
