package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Take(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Range(1, 9).Pipe(operators.Take(0)),
			Range(1, 9).Pipe(operators.Take(3)),
			Range(1, 3).Pipe(operators.Take(3)),
			Range(1, 1).Pipe(operators.Take(3)),
		},
		[][]interface{}{
			{Complete},
			{1, 2, 3, Complete},
			{1, 2, Complete},
			{Complete},
		},
	)
	subscribeN(
		t,
		[]Observable{
			Concat(Range(1, 9), Throw(errTest)).Pipe(operators.Take(0)),
			Concat(Range(1, 9), Throw(errTest)).Pipe(operators.Take(3)),
			Concat(Range(1, 3), Throw(errTest)).Pipe(operators.Take(3)),
			Concat(Range(1, 1), Throw(errTest)).Pipe(operators.Take(3)),
		},
		[][]interface{}{
			{Complete},
			{1, 2, 3, Complete},
			{1, 2, errTest},
			{errTest},
		},
	)
}
