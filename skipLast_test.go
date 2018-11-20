package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_SkipLast(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Range(1, 7).Pipe(operators.SkipLast(0)),
			Range(1, 7).Pipe(operators.SkipLast(3)),
			Range(1, 3).Pipe(operators.SkipLast(3)),
			Range(1, 1).Pipe(operators.SkipLast(3)),
		},
		1, 2, 3, 4, 5, 6, xComplete,
		1, 2, 3, xComplete,
		xComplete,
		xComplete,
	)

	subscribe(
		t,
		[]Observable{
			Concat(Range(1, 7), Throw(xErrTest)).Pipe(operators.SkipLast(0)),
			Concat(Range(1, 7), Throw(xErrTest)).Pipe(operators.SkipLast(3)),
			Concat(Range(1, 3), Throw(xErrTest)).Pipe(operators.SkipLast(3)),
			Concat(Range(1, 1), Throw(xErrTest)).Pipe(operators.SkipLast(3)),
		},
		1, 2, 3, 4, 5, 6, xErrTest,
		1, 2, 3, xErrTest,
		xErrTest,
		xErrTest,
	)
}
