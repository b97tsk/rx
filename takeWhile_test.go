package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_TakeWhile(t *testing.T) {
	takeLessThan5 := operators.TakeWhile(
		func(val interface{}, idx int) bool {
			return val.(int) < 5
		},
	)
	subscribe(
		t,
		[]Observable{
			Just(1, 2, 3, 4, 5, 4, 3, 2, 1).Pipe(takeLessThan5),
			Concat(Range(1, 9), Throw(xErrTest)).Pipe(takeLessThan5),
			Concat(Range(1, 5), Throw(xErrTest)).Pipe(takeLessThan5),
		},
		1, 2, 3, 4, xComplete,
		1, 2, 3, 4, xComplete,
		1, 2, 3, 4, xErrTest,
	)
}
