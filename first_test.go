package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_First(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Empty().Pipe(operators.First()),
			Throw(xErrTest).Pipe(operators.First()),
			Just("A").Pipe(operators.First()),
			Just("A", "B").Pipe(operators.First()),
			Concat(Just("A"), Throw(xErrTest)).Pipe(operators.First()),
			Concat(Just("A", "B"), Throw(xErrTest)).Pipe(operators.First()),
		},
		ErrEmpty,
		xErrTest,
		"A", xComplete,
		"A", xComplete,
		"A", xComplete,
		"A", xComplete,
	)
}
