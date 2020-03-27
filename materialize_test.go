package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Materialize(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Empty().Pipe(operators.Materialize(), operators.Count()),
			Throw(errTest).Pipe(operators.Materialize(), operators.Count()),
			Just("A", "B", "C").Pipe(operators.Materialize(), operators.Count()),
			Concat(Just("A", "B", "C"), Throw(errTest)).Pipe(operators.Materialize(), operators.Count()),
		},
		[][]interface{}{
			{1, Complete},
			{1, Complete},
			{4, Complete},
			{4, Complete},
		},
	)
}
