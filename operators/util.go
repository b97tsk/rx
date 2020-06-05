package operators

import (
	"github.com/b97tsk/rx"
)

func projectToObservable(val interface{}, idx int) (rx.Observable, error) {
	if obs, ok := val.(rx.Observable); ok {
		return obs, nil
	}
	return nil, rx.ErrNotObservable
}
