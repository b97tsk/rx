package operators

import (
	"github.com/b97tsk/rx"
)

func noop(source rx.Observable) rx.Observable {
	return source
}

func projectToObservable(val interface{}, idx int) rx.Observable {
	if obs, ok := val.(rx.Observable); ok {
		return obs
	}

	return rx.Throw(rx.ErrNotObservable)
}
