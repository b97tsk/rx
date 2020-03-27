package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Last(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Empty().Pipe(operators.Last()),
			Throw(errTest).Pipe(operators.Last()),
			Just("A").Pipe(operators.Last()),
			Just("A", "B").Pipe(operators.Last()),
			Concat(Just("A"), Throw(errTest)).Pipe(operators.Last()),
			Concat(Just("A", "B"), Throw(errTest)).Pipe(operators.Last()),
		},
		[][]interface{}{
			{ErrEmpty},
			{errTest},
			{"A", Complete},
			{"B", Complete},
			{errTest},
			{errTest},
		},
	)
}
