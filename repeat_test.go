package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Repeat(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Range(1, 4).Pipe(operators.Repeat(0)),
			Range(1, 4).Pipe(operators.Repeat(1)),
			Range(1, 4).Pipe(operators.Repeat(2)),
			Concat(Range(1, 4), Throw(xErrTest)).Pipe(operators.Repeat(0)),
			Concat(Range(1, 4), Throw(xErrTest)).Pipe(operators.Repeat(1)),
			Concat(Range(1, 4), Throw(xErrTest)).Pipe(operators.Repeat(2)),
		},
		xComplete,
		1, 2, 3, xComplete,
		1, 2, 3, 1, 2, 3, xComplete,
		xComplete,
		1, 2, 3, xErrTest,
		1, 2, 3, xErrTest,
	)
}
