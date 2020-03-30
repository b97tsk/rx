package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_IgnoreElements(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Empty().Pipe(operators.IgnoreElements()),
			Just("A", "B", "C").Pipe(operators.IgnoreElements()),
			Concat(Just("A", "B", "C"), Throw(errTest)).Pipe(operators.IgnoreElements()),
		},
		[][]interface{}{
			{Complete},
			{Complete},
			{errTest},
		},
	)
}
