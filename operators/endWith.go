package operators

import (
	"github.com/b97tsk/rx"
)

// EndWith creates an Observable that emits the items you specify as arguments
// after it finishes emitting items emitted by the source Observable.
func EndWith(values ...interface{}) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		if len(values) == 0 {
			return source
		}
		return rx.Concat(source, rx.FromSlice(values))
	}
}
