package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Find(t *testing.T) {
	findFive := operators.Find(
		func(val interface{}, idx int) bool {
			return val.(int) == 5
		},
	)
	subscribe(
		t,
		[]Observable{
			Just(1, 2, 3, 4, 5, 4, 3, 2, 1).Pipe(findFive),
			Range(1, 9).Pipe(findFive),
			Range(1, 5).Pipe(findFive),
			Concat(Range(1, 9), Throw(xErrTest)).Pipe(findFive),
			Concat(Range(1, 5), Throw(xErrTest)).Pipe(findFive),
		},
		5, xComplete,
		5, xComplete,
		xComplete,
		5, xComplete,
		xErrTest,
	)
}
