package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_FindIndex(t *testing.T) {
	findIndex := operators.FindIndex(
		func(val interface{}, idx int) bool {
			return val == "D"
		},
	)
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E").Pipe(findIndex),
			Just("A", "B", "C").Pipe(findIndex),
			Concat(Just("A", "B", "C", "D", "E"), Throw(xErrTest)).Pipe(findIndex),
			Concat(Just("A", "B", "C"), Throw(xErrTest)).Pipe(findIndex),
		},
		3, xComplete,
		xComplete,
		3, xComplete,
		xErrTest,
	)
}
