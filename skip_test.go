package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Skip(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Range(1, 7).Pipe(operators.Skip(0)),
			Range(1, 7).Pipe(operators.Skip(3)),
			Range(1, 3).Pipe(operators.Skip(3)),
			Range(1, 1).Pipe(operators.Skip(3)),
		},
		1, 2, 3, 4, 5, 6, xComplete,
		4, 5, 6, xComplete,
		xComplete,
		xComplete,
	)

	subscribe(
		t,
		[]Observable{
			Concat(Range(1, 7), Throw(xErrTest)).Pipe(operators.Skip(0)),
			Concat(Range(1, 7), Throw(xErrTest)).Pipe(operators.Skip(3)),
			Concat(Range(1, 3), Throw(xErrTest)).Pipe(operators.Skip(3)),
			Concat(Range(1, 1), Throw(xErrTest)).Pipe(operators.Skip(3)),
		},
		1, 2, 3, 4, 5, 6, xErrTest,
		4, 5, 6, xErrTest,
		xErrTest,
		xErrTest,
	)
}
