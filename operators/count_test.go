package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Count(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Empty().Pipe(operators.Count()),
			Range(1, 9).Pipe(operators.Count()),
			Concat(Range(1, 9), Throw(errTest)).Pipe(operators.Count()),
		},
		[][]interface{}{
			{0, Complete},
			{8, Complete},
			{errTest},
		},
	)
}
