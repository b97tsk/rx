package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_TakeLast(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Range(1, 9).Pipe(operators.TakeLast(0)),
			Range(1, 9).Pipe(operators.TakeLast(3)),
			Range(1, 3).Pipe(operators.TakeLast(3)),
			Range(1, 1).Pipe(operators.TakeLast(3)),
		},
		xComplete,
		6, 7, 8, xComplete,
		1, 2, xComplete,
		xComplete,
	)
	subscribe(
		t,
		[]Observable{
			Concat(Range(1, 9), Throw(xErrTest)).Pipe(operators.TakeLast(0)),
			Concat(Range(1, 9), Throw(xErrTest)).Pipe(operators.TakeLast(3)),
			Concat(Range(1, 3), Throw(xErrTest)).Pipe(operators.TakeLast(3)),
			Concat(Range(1, 1), Throw(xErrTest)).Pipe(operators.TakeLast(3)),
		},
		xComplete,
		xErrTest,
		xErrTest,
		xErrTest,
	)
}
