package operators

import (
	"github.com/b97tsk/rx"
)

// StartWith emits the items you specify as arguments before it begins to
// mirrors the source.
func StartWith(values ...interface{}) rx.Operator {
	if len(values) == 0 {
		return noop
	}

	return func(source rx.Observable) rx.Observable {
		return rx.Concat(rx.FromSlice(values), source)
	}
}
