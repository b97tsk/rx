package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_RepeatWhen(t *testing.T) {
	var (
		repeatOnce  = operators.Take(0)
		repeatTwice = operators.Take(1)
	)

	subscribe(
		t,
		[]Observable{
			Range(1, 4).Pipe(operators.RepeatWhen(repeatOnce)),
			Range(1, 4).Pipe(operators.RepeatWhen(repeatTwice)),
			Concat(Range(1, 4), Throw(xErrTest)).Pipe(operators.RepeatWhen(repeatOnce)),
			Concat(Range(1, 4), Throw(xErrTest)).Pipe(operators.RepeatWhen(repeatTwice)),
		},
		1, 2, 3, xComplete,
		1, 2, 3, 1, 2, 3, xComplete,
		1, 2, 3, xErrTest,
		1, 2, 3, xErrTest,
	)
}
