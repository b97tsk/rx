package operators_test

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
	subscribeN(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E").Pipe(findIndex),
			Just("A", "B", "C").Pipe(findIndex),
			Concat(Just("A", "B", "C", "D", "E"), Throw(errTest)).Pipe(findIndex),
			Concat(Just("A", "B", "C"), Throw(errTest)).Pipe(findIndex),
		},
		[][]interface{}{
			{3, Complete},
			{Complete},
			{3, Complete},
			{errTest},
		},
	)
}
