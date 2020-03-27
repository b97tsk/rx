package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_StartWith(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Just("D", "E").Pipe(operators.StartWith("A", "B", "C")),
			Empty().Pipe(operators.StartWith("A", "B", "C")),
			Throw(errTest).Pipe(operators.StartWith("A", "B", "C")),
			Throw(errTest).Pipe(operators.StartWith()),
		},
		[][]interface{}{
			{"A", "B", "C", "D", "E", Complete},
			{"A", "B", "C", Complete},
			{"A", "B", "C", errTest},
			{errTest},
		},
	)
}
