package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_EndWith(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C").Pipe(operators.EndWith("D", "E")),
			Empty().Pipe(operators.EndWith("D", "E")),
			Throw(xErrTest).Pipe(operators.EndWith("D", "E")),
			Throw(xErrTest).Pipe(operators.EndWith()),
		},
		"A", "B", "C", "D", "E", xComplete,
		"D", "E", xComplete,
		xErrTest,
		xErrTest,
	)
}
