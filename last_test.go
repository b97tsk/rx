package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Last(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Empty().Pipe(operators.Last()),
			Throw(xErrTest).Pipe(operators.Last()),
			Just("A").Pipe(operators.Last()),
			Just("A", "B").Pipe(operators.Last()),
			Concat(Just("A"), Throw(xErrTest)).Pipe(operators.Last()),
			Concat(Just("A", "B"), Throw(xErrTest)).Pipe(operators.Last()),
		},
		ErrEmpty,
		xErrTest,
		"A", xComplete,
		"B", xComplete,
		xErrTest,
		xErrTest,
	)
}
