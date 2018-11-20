package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Materialize(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Empty().Pipe(operators.Materialize(), operators.Count()),
			Throw(xErrTest).Pipe(operators.Materialize(), operators.Count()),
			Just("A", "B", "C").Pipe(operators.Materialize(), operators.Count()),
			Concat(Just("A", "B", "C"), Throw(xErrTest)).Pipe(operators.Materialize(), operators.Count()),
		},
		1, xComplete,
		1, xComplete,
		4, xComplete,
		4, xComplete,
	)
}
