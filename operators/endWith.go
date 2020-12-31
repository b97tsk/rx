package operators

import (
	"github.com/b97tsk/rx"
)

// EndWith mirrors the source, emits the items you specify as arguments when
// the source completes.
func EndWith(values ...interface{}) rx.Operator {
	if len(values) == 0 {
		return noop
	}

	return func(source rx.Observable) rx.Observable {
		return rx.Concat(source, rx.FromSlice(values))
	}
}
