package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Skip(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Range(1, 7).Pipe(operators.Skip(0)),
			Range(1, 7).Pipe(operators.Skip(3)),
			Range(1, 3).Pipe(operators.Skip(3)),
			Range(1, 1).Pipe(operators.Skip(3)),
		},
		[][]interface{}{
			{1, 2, 3, 4, 5, 6, Complete},
			{4, 5, 6, Complete},
			{Complete},
			{Complete},
		},
	)

	subscribeN(
		t,
		[]Observable{
			Concat(Range(1, 7), Throw(errTest)).Pipe(operators.Skip(0)),
			Concat(Range(1, 7), Throw(errTest)).Pipe(operators.Skip(3)),
			Concat(Range(1, 3), Throw(errTest)).Pipe(operators.Skip(3)),
			Concat(Range(1, 1), Throw(errTest)).Pipe(operators.Skip(3)),
		},
		[][]interface{}{
			{1, 2, 3, 4, 5, 6, errTest},
			{4, 5, 6, errTest},
			{errTest},
			{errTest},
		},
	)
}
