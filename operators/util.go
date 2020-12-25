package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

func noop(source rx.Observable) rx.Observable {
	return source
}

func projectToObservable(val interface{}, idx int) rx.Observable {
	switch obs := val.(type) {
	case rx.Observable:
		return obs
	case func(context.Context, rx.Observer):
		return obs
	default:
		return rx.Throw(rx.ErrNotObservable)
	}
}
