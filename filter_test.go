package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Filter(t *testing.T) {
	filterLessThan5 := operators.Filter(
		func(val interface{}, idx int) bool {
			return val.(int) < 5
		},
	)
	subscribe(
		t,
		[]Observable{
			Just(1, 2, 3, 4, 5, 4, 3, 2, 1).Pipe(filterLessThan5),
			Range(1, 9).Pipe(filterLessThan5),
			Concat(Range(1, 9), Throw(xErrTest)).Pipe(filterLessThan5),
		},
		1, 2, 3, 4, 4, 3, 2, 1, xComplete,
		1, 2, 3, 4, xComplete,
		1, 2, 3, 4, xErrTest,
	)
}
