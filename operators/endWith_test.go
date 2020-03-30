package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_EndWith(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Just("A", "B", "C").Pipe(operators.EndWith("D", "E")),
			Empty().Pipe(operators.EndWith("D", "E")),
			Throw(errTest).Pipe(operators.EndWith("D", "E")),
			Throw(errTest).Pipe(operators.EndWith()),
		},
		[][]interface{}{
			{"A", "B", "C", "D", "E", Complete},
			{"D", "E", Complete},
			{errTest},
			{errTest},
		},
	)
}
