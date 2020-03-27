package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Single(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Just("A", "B").Pipe(operators.Single()),
			Just("A").Pipe(operators.Single()),
			Empty().Pipe(operators.Single()),
			Throw(errTest).Pipe(operators.Single()),
		},
		[][]interface{}{
			{ErrNotSingle},
			{"A", Complete},
			{ErrEmpty},
			{errTest},
		},
	)
}
