package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_StartWith(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Just("D", "E").Pipe(operators.StartWith("A", "B", "C")),
			Empty().Pipe(operators.StartWith("A", "B", "C")),
			Throw(xErrTest).Pipe(operators.StartWith("A", "B", "C")),
			Throw(xErrTest).Pipe(operators.StartWith()),
		},
		"A", "B", "C", "D", "E", xComplete,
		"A", "B", "C", xComplete,
		"A", "B", "C", xErrTest,
		xErrTest,
	)
}
