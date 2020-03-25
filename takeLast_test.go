package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_TakeLast(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Range(1, 9).Pipe(operators.TakeLast(0)),
			Range(1, 9).Pipe(operators.TakeLast(3)),
			Range(1, 3).Pipe(operators.TakeLast(3)),
			Range(1, 1).Pipe(operators.TakeLast(3)),
		},
		Complete,
		6, 7, 8, Complete,
		1, 2, Complete,
		Complete,
	)
	subscribe(
		t,
		[]Observable{
			Concat(Range(1, 9), Throw(errTest)).Pipe(operators.TakeLast(0)),
			Concat(Range(1, 9), Throw(errTest)).Pipe(operators.TakeLast(3)),
			Concat(Range(1, 3), Throw(errTest)).Pipe(operators.TakeLast(3)),
			Concat(Range(1, 1), Throw(errTest)).Pipe(operators.TakeLast(3)),
		},
		Complete,
		errTest,
		errTest,
		errTest,
	)
}
