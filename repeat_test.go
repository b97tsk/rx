package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Repeat(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Range(1, 4).Pipe(operators.Repeat(0)),
			Range(1, 4).Pipe(operators.Repeat(1)),
			Range(1, 4).Pipe(operators.Repeat(2)),
			Concat(Range(1, 4), Throw(errTest)).Pipe(operators.Repeat(0)),
			Concat(Range(1, 4), Throw(errTest)).Pipe(operators.Repeat(1)),
			Concat(Range(1, 4), Throw(errTest)).Pipe(operators.Repeat(2)),
		},
		[][]interface{}{
			{Complete},
			{1, 2, 3, Complete},
			{1, 2, 3, 1, 2, 3, Complete},
			{Complete},
			{1, 2, 3, errTest},
			{1, 2, 3, errTest},
		},
	)
}
