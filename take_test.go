package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Take(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Range(1, 9).Pipe(operators.Take(0)),
			Range(1, 9).Pipe(operators.Take(3)),
			Range(1, 3).Pipe(operators.Take(3)),
			Range(1, 1).Pipe(operators.Take(3)),
		},
		xComplete,
		1, 2, 3, xComplete,
		1, 2, xComplete,
		xComplete,
	)
	subscribe(
		t,
		[]Observable{
			Concat(Range(1, 9), Throw(xErrTest)).Pipe(operators.Take(0)),
			Concat(Range(1, 9), Throw(xErrTest)).Pipe(operators.Take(3)),
			Concat(Range(1, 3), Throw(xErrTest)).Pipe(operators.Take(3)),
			Concat(Range(1, 1), Throw(xErrTest)).Pipe(operators.Take(3)),
		},
		xComplete,
		1, 2, 3, xComplete,
		1, 2, xErrTest,
		xErrTest,
	)
}
